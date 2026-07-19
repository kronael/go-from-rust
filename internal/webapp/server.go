package webapp

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultPlaygroundURL = "https://go.dev/_/compile"
	defaultUserAgent     = "go-from-rust/1 (+https://github.com/kronael/go-from-rust)"

	maxSourceBytes    = 128 * 1024
	maxRequestBytes   = 512 * 1024
	maxUpstreamBytes  = 1024 * 1024
	cacheCapacity     = 256
	cacheTTL          = 10 * time.Minute
	maxConcurrentRuns = 8
	runBudgetPerMin   = 120
	upstreamTimeout   = 12 * time.Second
)

// Config configures a Server. Zero values fall back to production defaults;
// tests override HTTPClient and Now to control timeouts and cache/budget expiry.
type Config struct {
	PlaygroundURL string
	HTTPClient    *http.Client
	Now           func() time.Time
}

type cacheEntry struct {
	value     []byte
	expiresAt time.Time
}

// Server implements the guarded web API: the lesson catalog, gofmt, and a
// rate-limited, cached proxy to the Go Playground compile API.
type Server struct {
	lessons       []Lesson
	webFS         fs.FS
	playgroundURL string
	httpClient    *http.Client
	now           func() time.Time
	mux           *http.ServeMux

	mu          sync.Mutex
	cache       map[string]cacheEntry
	budgetStart time.Time
	budgetUsed  int
	sem         chan struct{}
}

// NewServer builds a Server serving the given lesson catalog and web assets.
func NewServer(lessons []Lesson, webFS fs.FS, cfg Config) *Server {
	playgroundURL := cfg.PlaygroundURL
	if playgroundURL == "" {
		playgroundURL = defaultPlaygroundURL
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: upstreamTimeout}
	}
	now := cfg.Now
	if now == nil {
		now = time.Now
	}

	s := &Server{
		lessons:       lessons,
		webFS:         webFS,
		playgroundURL: playgroundURL,
		httpClient:    httpClient,
		now:           now,
		cache:         make(map[string]cacheEntry),
		sem:           make(chan struct{}, maxConcurrentRuns),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", s.handleIndex)
	mux.HandleFunc("/static/app.js", s.handleStatic("app.js", "application/javascript; charset=utf-8"))
	mux.HandleFunc("/static/styles.css", s.handleStatic("styles.css", "text/css; charset=utf-8"))
	mux.HandleFunc("/static/go-logo-white.svg", s.handleStatic("go-logo-white.svg", "image/svg+xml"))
	mux.HandleFunc("/static/gopher.png", s.handleStatic("gopher.png", "image/png"))
	mux.HandleFunc("/api/lessons", s.handleLessons)
	mux.HandleFunc("/api/format", s.handleFormat)
	mux.HandleFunc("/api/run", s.handleRun)
	mux.HandleFunc("/health", s.handleHealth)
	s.mux = mux

	return s
}

// Handler returns all routes with browser security headers.
func (s *Server) Handler() http.Handler {
	return s.securityHeaders(s.mux)
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; base-uri 'none'; connect-src 'self'; form-action 'none'; frame-ancestors 'none'; img-src 'self' data:; script-src 'self'; style-src 'self'")
		w.Header().Set("Permissions-Policy", "camera=(), geolocation=(), microphone=()")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

// --- error envelope ---

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{Error: errorBody{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// --- static routes ---

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is supported")
		return
	}
	s.serveEmbedded(w, "index.html", "text/html; charset=utf-8", "no-cache")
}

func (s *Server) handleStatic(name, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is supported")
			return
		}
		s.serveEmbedded(w, name, contentType, "public, max-age=3600")
	}
}

func (s *Server) serveEmbedded(w http.ResponseWriter, name, contentType, cacheControl string) {
	data, err := fs.ReadFile(s.webFS, name)
	if err != nil {
		http.NotFound(w, nil)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", cacheControl)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// --- health ---

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is supported")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// --- lessons ---

func (s *Server) handleLessons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is supported")
		return
	}
	writeJSON(w, http.StatusOK, map[string][]Lesson{"lessons": s.lessons})
}

// --- format ---

type bodyRequest struct {
	Body string `json:"body"`
}

type bodyResponse struct {
	Body string `json:"body"`
}

func decodeBodyRequest(r *http.Request) (bodyRequest, int, string, string) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxRequestBytes+1))
	if err != nil || len(raw) > maxRequestBytes {
		return bodyRequest{}, http.StatusBadRequest, "invalid_request", "request body could not be read"
	}
	var payload struct {
		Body *string `json:"body"`
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		return bodyRequest{}, http.StatusBadRequest, "invalid_request", "request body must be JSON with a \"body\" field"
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return bodyRequest{}, http.StatusBadRequest, "invalid_request", "request body must contain one JSON object"
	}
	if payload.Body == nil {
		return bodyRequest{}, http.StatusBadRequest, "invalid_request", "request body must include a string \"body\" field"
	}
	req := bodyRequest{Body: *payload.Body}
	if len(req.Body) > maxSourceBytes {
		return bodyRequest{}, http.StatusRequestEntityTooLarge, "source_too_large", "source exceeds 128KiB"
	}
	return req, 0, "", ""
}

func (s *Server) handleFormat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST is supported")
		return
	}
	req, status, code, msg := decodeBodyRequest(r)
	if status != 0 {
		writeError(w, status, code, msg)
		return
	}
	formatted, err := format.Source([]byte(req.Body))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "format_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, bodyResponse{Body: string(formatted)})
}

// --- run ---

func hashSource(body string) string {
	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
}

func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST is supported")
		return
	}
	req, status, code, msg := decodeBodyRequest(r)
	if status != 0 {
		writeError(w, status, code, msg)
		return
	}

	key := hashSource(req.Body)
	if value, ok := s.cacheGet(key); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(value)
		return
	}

	select {
	case s.sem <- struct{}{}:
	default:
		w.Header().Set("Retry-After", "1")
		writeError(w, http.StatusServiceUnavailable, "runner_busy", "too many concurrent runs, try again shortly")
		return
	}
	defer func() { <-s.sem }()

	if !s.reserveBudget() {
		w.Header().Set("Retry-After", "60")
		writeError(w, http.StatusTooManyRequests, "run_budget_exhausted", "run budget exhausted, try again later")
		return
	}

	value, status, code, msg := s.runUpstream(r.Context(), req.Body)
	if status != http.StatusOK {
		writeError(w, status, code, msg)
		return
	}

	s.cachePut(key, value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(value)
}

func (s *Server) runUpstream(ctx context.Context, source string) ([]byte, int, string, string) {
	form := url.Values{}
	form.Set("version", "2")
	form.Set("body", source)
	form.Set("withVet", "true")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.playgroundURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, http.StatusBadGateway, "upstream_failed", "failed to build upstream request"
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, http.StatusGatewayTimeout, "upstream_timeout", "upstream request timed out"
		}
		return nil, http.StatusBadGateway, "upstream_failed", "upstream request failed"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, http.StatusBadGateway, "upstream_failed", fmt.Sprintf("upstream returned status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxUpstreamBytes+1))
	if err != nil || len(body) > maxUpstreamBytes {
		return nil, http.StatusBadGateway, "upstream_failed", "upstream response too large or unreadable"
	}
	if !json.Valid(body) {
		return nil, http.StatusBadGateway, "upstream_failed", "upstream response was not valid JSON"
	}
	return body, http.StatusOK, "", ""
}

// --- cache ---

func (s *Server) cacheGet(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.cache[key]
	if !ok {
		return nil, false
	}
	if s.now().After(entry.expiresAt) {
		delete(s.cache, key)
		return nil, false
	}
	return entry.value, true
}

func (s *Server) cachePut(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	for cachedKey, entry := range s.cache {
		if now.After(entry.expiresAt) {
			delete(s.cache, cachedKey)
		}
	}
	if _, exists := s.cache[key]; !exists && len(s.cache) >= cacheCapacity {
		clear(s.cache)
	}
	s.cache[key] = cacheEntry{value: value, expiresAt: now.Add(cacheTTL)}
}

// --- budget ---

func (s *Server) reserveBudget() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()
	if s.budgetStart.IsZero() || !now.Before(s.budgetStart.Add(time.Minute)) {
		s.budgetStart = now
		s.budgetUsed = 0
	}
	if s.budgetUsed >= runBudgetPerMin {
		return false
	}
	s.budgetUsed++
	return true
}

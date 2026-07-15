package webapp

import (
	"bytes"
	"container/list"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultPlaygroundURL is the official Go Playground compile endpoint.
	DefaultPlaygroundURL = "https://go.dev/_/compile"
	// DefaultUserAgent identifies this proxy to the Playground.
	DefaultUserAgent = "go-from-rust/1 (+https://github.com/kronael/go-from-rust)"

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
	UserAgent     string
	Logger        *log.Logger
	HTTPClient    *http.Client
	Now           func() time.Time
}

type cacheEntry struct {
	key       string
	value     []byte
	expiresAt time.Time
}

// Server implements the guarded web API: the lesson catalog, gofmt, and a
// rate-limited, cached proxy to the Go Playground compile API.
type Server struct {
	lessons       []Lesson
	webFS         fs.FS
	playgroundURL string
	userAgent     string
	httpClient    *http.Client
	logger        *log.Logger
	now           func() time.Time
	mux           *http.ServeMux

	mu          sync.Mutex
	cacheList   *list.List
	cacheMap    map[string]*list.Element
	budgetTimes []time.Time
	sem         chan struct{}
}

// NewServer builds a Server serving the given lesson catalog and embedded web
// filesystem (containing index.html, app.js, styles.css at its root).
func NewServer(lessons []Lesson, webFS fs.FS, cfg Config) *Server {
	playgroundURL := cfg.PlaygroundURL
	if playgroundURL == "" {
		playgroundURL = DefaultPlaygroundURL
	}
	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = DefaultUserAgent
	}
	logger := cfg.Logger
	if logger == nil {
		logger = log.Default()
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
		userAgent:     userAgent,
		httpClient:    httpClient,
		logger:        logger,
		now:           now,
		cacheList:     list.New(),
		cacheMap:      make(map[string]*list.Element),
		sem:           make(chan struct{}, maxConcurrentRuns),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", s.handleIndex)
	mux.HandleFunc("/static/app.js", s.handleStaticJS)
	mux.HandleFunc("/static/styles.css", s.handleStaticCSS)
	mux.HandleFunc("/api/lessons", s.handleLessons)
	mux.HandleFunc("/api/format", s.handleFormat)
	mux.HandleFunc("/api/run", s.handleRun)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)
	s.mux = mux

	return s
}

// Handler returns the http.Handler serving all routes, wrapped with security
// headers and request logging.
func (s *Server) Handler() http.Handler {
	return s.logRequests(s.securityHeaders(s.mux))
}

func (s *Server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := s.now()
		reqID := requestID(r)
		w.Header().Set("X-Request-ID", reqID)
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		s.logger.Printf("request_id=%s method=%s path=%q status=%d duration=%s",
			reqID, r.Method, r.URL.Path, rec.status, s.now().Sub(start))
	})
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

type statusRecorder struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.wrote {
		return
	}
	r.wrote = true
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func requestID(r *http.Request) string {
	id := r.Header.Get("X-Request-ID")
	if validRequestID(id) {
		return id
	}
	return newRequestID()
}

func validRequestID(id string) bool {
	if len(id) == 0 || len(id) > 64 {
		return false
	}
	for _, char := range id {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || strings.ContainsRune("._-", char) {
			continue
		}
		return false
	}
	return true
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
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

func (s *Server) handleStaticJS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is supported")
		return
	}
	s.serveEmbedded(w, "app.js", "application/javascript; charset=utf-8", "public, max-age=3600")
}

func (s *Server) handleStaticCSS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is supported")
		return
	}
	s.serveEmbedded(w, "styles.css", "text/css; charset=utf-8", "public, max-age=3600")
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

// --- health/ready ---

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is supported")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is supported")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ready": true})
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
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || isTimeoutErr(err) {
			return nil, http.StatusGatewayTimeout, "upstream_timeout", "upstream request timed out"
		}
		return nil, http.StatusBadGateway, "upstream_failed", "upstream request failed"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, http.StatusBadGateway, "upstream_failed", fmt.Sprintf("upstream returned status %d", resp.StatusCode)
	}
	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return nil, http.StatusBadGateway, "upstream_failed", "upstream response was not JSON"
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

func isTimeoutErr(err error) bool {
	type timeouter interface{ Timeout() bool }
	var t timeouter
	if errors.As(err, &t) {
		return t.Timeout()
	}
	return strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Client.Timeout")
}

// --- cache ---

func (s *Server) cacheGet(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	elem, ok := s.cacheMap[key]
	if !ok {
		return nil, false
	}
	entry := elem.Value.(*cacheEntry)
	if s.now().After(entry.expiresAt) {
		s.cacheList.Remove(elem)
		delete(s.cacheMap, key)
		return nil, false
	}
	s.cacheList.MoveToFront(elem)
	return entry.value, true
}

func (s *Server) cachePut(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.purgeExpiredLocked()

	if elem, ok := s.cacheMap[key]; ok {
		entry := elem.Value.(*cacheEntry)
		entry.value = value
		entry.expiresAt = s.now().Add(cacheTTL)
		s.cacheList.MoveToFront(elem)
		return
	}

	if s.cacheList.Len() >= cacheCapacity {
		back := s.cacheList.Back()
		if back != nil {
			be := back.Value.(*cacheEntry)
			delete(s.cacheMap, be.key)
			s.cacheList.Remove(back)
		}
	}

	entry := &cacheEntry{key: key, value: value, expiresAt: s.now().Add(cacheTTL)}
	elem := s.cacheList.PushFront(entry)
	s.cacheMap[key] = elem
}

func (s *Server) purgeExpiredLocked() {
	var next *list.Element
	for e := s.cacheList.Front(); e != nil; e = next {
		next = e.Next()
		entry := e.Value.(*cacheEntry)
		if s.now().After(entry.expiresAt) {
			s.cacheList.Remove(e)
			delete(s.cacheMap, entry.key)
		}
	}
}

// --- budget ---

func (s *Server) reserveBudget() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := s.now().Add(-time.Minute)
	i := 0
	for i < len(s.budgetTimes) && s.budgetTimes[i].Before(cutoff) {
		i++
	}
	s.budgetTimes = s.budgetTimes[i:]

	if len(s.budgetTimes) >= runBudgetPerMin {
		return false
	}
	s.budgetTimes = append(s.budgetTimes, s.now())
	return true
}

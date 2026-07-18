package webapp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/fstest"
	"time"
)

func testWebFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html>tour</html>")},
		"app.js":     &fstest.MapFile{Data: []byte("console.log('tour')")},
		"styles.css": &fstest.MapFile{Data: []byte("body{}")},
	}
}

func testLessons() []Lesson {
	return []Lesson{
		{ID: 1, Filename: "01_arrays_slices.go", Value: 5, Question: "Q1?", Source: "package main\n"},
		{ID: 2, Filename: "02_copy.go", Value: 5, Question: "Q2?", Source: "package main\n"},
	}
}

// fakePlayground is a configurable in-process stand-in for the Go Playground
// compile endpoint. Real network calls are never made in these tests.
type fakePlayground struct {
	delay    time.Duration
	status   int
	response string
	requests int32
}

func newFakePlayground() *fakePlayground {
	return &fakePlayground{status: http.StatusOK, response: `{"Events":[{"Message":"hi"}]}`}
}

func (f *fakePlayground) server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&f.requests, 1)
		if r.Method != http.MethodPost || r.ParseForm() != nil || r.Form.Get("version") != "2" || r.Form.Get("withVet") != "true" || r.Header.Get("User-Agent") != defaultUserAgent {
			http.Error(w, "bad playground request", http.StatusBadRequest)
			return
		}

		if f.delay > 0 {
			select {
			case <-time.After(f.delay):
			case <-r.Context().Done():
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(f.status)
		_, _ = w.Write([]byte(f.response))
	}))
}

func newTestServer(t *testing.T, playgroundURL string, now func() time.Time) *Server {
	t.Helper()
	cfg := Config{
		PlaygroundURL: playgroundURL,
		HTTPClient:    &http.Client{Timeout: 2 * time.Second},
		Now:           now,
	}
	return NewServer(testLessons(), testWebFS(), cfg)
}

func decodeError(t *testing.T, body io.Reader) errorEnvelope {
	t.Helper()
	var env errorEnvelope
	if err := json.NewDecoder(body).Decode(&env); err != nil {
		t.Fatalf("decode error envelope: %v", err)
	}
	return env
}

func TestLessonsEndpoint(t *testing.T) {
	srv := httptest.NewServer(newTestServer(t, "http://unused.invalid", time.Now).Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/lessons")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var out struct {
		Lessons []Lesson `json:"lessons"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if len(out.Lessons) != 2 {
		t.Fatalf("got %d lessons, want 2", len(out.Lessons))
	}
}

func TestLessonsMethodNotAllowed(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/lessons", "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", resp.StatusCode)
	}
	env := decodeError(t, resp.Body)
	if env.Error.Code != "method_not_allowed" {
		t.Fatalf("code = %q, want method_not_allowed", env.Error.Code)
	}
}

func TestFormatSuccess(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\nfunc main(){}\n"}`
	resp, err := http.Post(srv.URL+"/api/format", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var out bodyResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.Body, "func main()") {
		t.Fatalf("unexpected formatted body: %q", out.Body)
	}
}

func TestFormatInvalidJSON(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/format", "application/json", strings.NewReader("not json"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
	if env := decodeError(t, resp.Body); env.Error.Code != "invalid_request" {
		t.Fatalf("code = %q, want invalid_request", env.Error.Code)
	}
}

func TestFormatRequiresExactlyOneBodyField(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	for _, body := range []string{`{}`, `{"body":null}`, `{"body":"package main","extra":true}`, `{"body":"package main"}{}`} {
		resp, err := http.Post(srv.URL+"/api/format", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("body %q: status = %d, want 400", body, resp.StatusCode)
		}
	}
}

func TestFormatUnformattable(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\nfunc main( {\n"}`
	resp, err := http.Post(srv.URL+"/api/format", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", resp.StatusCode)
	}
	if env := decodeError(t, resp.Body); env.Error.Code != "format_failed" {
		t.Fatalf("code = %q, want format_failed", env.Error.Code)
	}
}

func TestFormatSourceTooLarge(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	huge := strings.Repeat("a", maxSourceBytes+1)
	payload, err := json.Marshal(bodyRequest{Body: huge})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(srv.URL+"/api/format", "application/json", strings.NewReader(string(payload)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", resp.StatusCode)
	}
	if env := decodeError(t, resp.Body); env.Error.Code != "source_too_large" {
		t.Fatalf("code = %q, want source_too_large", env.Error.Code)
	}
}

func TestRunSuccessAndCacheHit(t *testing.T) {
	fp := newFakePlayground()
	fps := fp.server()
	defer fps.Close()

	s := newTestServer(t, fps.URL, time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\nfunc main(){}\n"}`

	resp1, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp1.StatusCode)
	}
	b1, _ := io.ReadAll(resp1.Body)

	resp2, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp2.StatusCode)
	}
	b2, _ := io.ReadAll(resp2.Body)

	if string(b1) != string(b2) {
		t.Fatalf("cached response mismatch: %q vs %q", b1, b2)
	}
	if atomic.LoadInt32(&fp.requests) != 1 {
		t.Fatalf("upstream requests = %d, want 1 (second call should hit cache)", fp.requests)
	}
}

func TestRunInvalidJSON(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader("{"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", resp.StatusCode)
	}
}

func TestRunSourceTooLarge(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	huge := strings.Repeat("a", maxSourceBytes+1)
	payload, _ := json.Marshal(bodyRequest{Body: huge})
	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(string(payload)))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", resp.StatusCode)
	}
}

func TestRunMethodNotAllowed(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/run")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", resp.StatusCode)
	}
}

func TestRunUpstreamBadStatus(t *testing.T) {
	fp := newFakePlayground()
	fp.status = http.StatusInternalServerError
	fps := fp.server()
	defer fps.Close()

	s := newTestServer(t, fps.URL, time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\n"}`
	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", resp.StatusCode)
	}
	if env := decodeError(t, resp.Body); env.Error.Code != "upstream_failed" {
		t.Fatalf("code = %q, want upstream_failed", env.Error.Code)
	}
}

func TestRunUpstreamInvalidJSON(t *testing.T) {
	fp := newFakePlayground()
	fp.response = "not json"
	fps := fp.server()
	defer fps.Close()

	s := newTestServer(t, fps.URL, time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\n"}`
	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", resp.StatusCode)
	}
	if env := decodeError(t, resp.Body); env.Error.Code != "upstream_failed" {
		t.Fatalf("code = %q, want upstream_failed", env.Error.Code)
	}
}

func TestRunUpstreamAcceptsPlaygroundContentType(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(`{"Events":[]}`))
	}))
	defer upstream.Close()

	srv := httptest.NewServer(newTestServer(t, upstream.URL, time.Now).Handler())
	defer srv.Close()
	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(`{"body":"package main\n"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var payload map[string][]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if _, ok := payload["Events"]; !ok {
		t.Fatal("response omitted Events")
	}
}

func TestRunUpstreamOversized(t *testing.T) {
	fp := newFakePlayground()
	fp.response = `{"pad":"` + strings.Repeat("a", maxUpstreamBytes+1) + `"}`
	fps := fp.server()
	defer fps.Close()

	s := newTestServer(t, fps.URL, time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\n"}`
	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", resp.StatusCode)
	}
	if env := decodeError(t, resp.Body); env.Error.Code != "upstream_failed" {
		t.Fatalf("code = %q, want upstream_failed", env.Error.Code)
	}
}

func TestRunUpstreamTimeout(t *testing.T) {
	fp := newFakePlayground()
	fp.delay = 200 * time.Millisecond
	fps := fp.server()
	defer fps.Close()

	cfg := Config{
		PlaygroundURL: fps.URL,
		HTTPClient:    &http.Client{Timeout: 20 * time.Millisecond},
		Now:           time.Now,
	}
	s := NewServer(testLessons(), testWebFS(), cfg)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\n"}`
	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusGatewayTimeout {
		t.Fatalf("status = %d, want 504", resp.StatusCode)
	}
	if env := decodeError(t, resp.Body); env.Error.Code != "upstream_timeout" {
		t.Fatalf("code = %q, want upstream_timeout", env.Error.Code)
	}
}

func TestRunConcurrencyLimit(t *testing.T) {
	fp := newFakePlayground()
	release := make(chan struct{})
	fps := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&fp.requests, 1)
		<-release
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Events":[]}`))
	}))
	defer fps.Close()

	s := newTestServer(t, fps.URL, time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	var wg sync.WaitGroup
	statuses := make([]int, maxConcurrentRuns+1)
	for i := 0; i < maxConcurrentRuns+1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			body := `{"body":"package main\n// ` + strings.Repeat("x", i+1) + `\n"}`
			resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			statuses[i] = resp.StatusCode
		}(i)
	}

	// Give all goroutines time to reach the upstream (or be turned away) before releasing.
	deadline := time.After(2 * time.Second)
	for {
		if atomic.LoadInt32(&fp.requests) >= maxConcurrentRuns {
			break
		}
		select {
		case <-deadline:
			t.Fatal("timed out waiting for concurrent requests to reach upstream")
		case <-time.After(5 * time.Millisecond):
		}
	}
	close(release)
	wg.Wait()

	var busy, ok int
	for _, status := range statuses {
		switch status {
		case http.StatusOK:
			ok++
		case http.StatusServiceUnavailable:
			busy++
		}
	}
	if busy < 1 {
		t.Fatalf("expected at least one 503 runner_busy, got statuses %v", statuses)
	}
	if ok != maxConcurrentRuns {
		t.Fatalf("expected %d successful runs, got %d (statuses %v)", maxConcurrentRuns, ok, statuses)
	}
}

func TestRunBudgetExhausted(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)

	// Fill the minute budget directly instead of firing 120 requests.
	s.mu.Lock()
	s.budgetStart = s.now()
	s.budgetUsed = runBudgetPerMin
	s.mu.Unlock()

	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\n"}`
	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429", resp.StatusCode)
	}
	if got := resp.Header.Get("Retry-After"); got != "60" {
		t.Fatalf("Retry-After = %q, want 60", got)
	}
	if env := decodeError(t, resp.Body); env.Error.Code != "run_budget_exhausted" {
		t.Fatalf("code = %q, want run_budget_exhausted", env.Error.Code)
	}
}

func TestRunCacheHitBypassesBudget(t *testing.T) {
	fp := newFakePlayground()
	fps := fp.server()
	defer fps.Close()

	s := newTestServer(t, fps.URL, time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	body := `{"body":"package main\nfunc main(){}\n"}`
	resp, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("priming request status = %d, want 200", resp.StatusCode)
	}

	// Exhaust the budget after the entry is cached.
	s.mu.Lock()
	s.budgetStart = s.now()
	s.budgetUsed = runBudgetPerMin
	s.mu.Unlock()

	resp2, err := http.Post(srv.URL+"/api/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("cached request status = %d, want 200 (cache hits must bypass budget)", resp2.StatusCode)
	}
}

func TestCacheExpiryAndCapacity(t *testing.T) {
	current := time.Now()
	clock := func() time.Time { return current }
	s := newTestServer(t, "http://unused.invalid", clock)

	s.cachePut("key-a", []byte(`{"a":1}`))
	current = current.Add(cacheTTL + time.Second)
	if _, ok := s.cacheGet("key-a"); ok {
		t.Fatal("expected expired cache entry to be evicted")
	}

	current = time.Now()
	for i := 0; i < cacheCapacity; i++ {
		s.cachePut(keyN(i), []byte(`{}`))
	}
	newest := keyN(cacheCapacity)
	s.cachePut(newest, []byte(`{}`))
	if _, ok := s.cacheGet(keyN(0)); ok {
		t.Fatal("expected full cache to be cleared once over capacity")
	}
	if _, ok := s.cacheGet(newest); !ok {
		t.Fatal("expected newest entry to survive cache reset")
	}
}

func keyN(i int) string {
	return "key-" + string(rune('a'+i%26)) + string(rune('0'+i/26))
}

func TestHealth(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

func TestStaticRoutesAndUnknownPath(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	srv := httptest.NewServer(s.Handler())
	defer srv.Close()

	cases := []struct {
		path        string
		contentType string
		cache       string
	}{
		{"/", "text/html; charset=utf-8", "no-cache"},
		{"/static/app.js", "application/javascript; charset=utf-8", "public, max-age=3600"},
		{"/static/styles.css", "text/css; charset=utf-8", "public, max-age=3600"},
	}
	for _, c := range cases {
		resp, err := http.Get(srv.URL + c.path)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s status = %d, want 200", c.path, resp.StatusCode)
		}
		if got := resp.Header.Get("Content-Type"); got != c.contentType {
			t.Fatalf("%s content-type = %q, want %q", c.path, got, c.contentType)
		}
		if got := resp.Header.Get("Cache-Control"); got != c.cache {
			t.Fatalf("%s cache-control = %q, want %q", c.path, got, c.cache)
		}
	}

	resp, err := http.Get(srv.URL + "/nope")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unknown path status = %d, want 404", resp.StatusCode)
	}
}

func TestSecurityHeaders(t *testing.T) {
	s := newTestServer(t, "http://unused.invalid", time.Now)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	s.Handler().ServeHTTP(rec, req)

	for _, header := range []string{
		"Content-Security-Policy",
		"Permissions-Policy",
		"Referrer-Policy",
		"X-Content-Type-Options",
		"X-Frame-Options",
	} {
		if rec.Header().Get(header) == "" {
			t.Errorf("%s header is empty", header)
		}
	}
}

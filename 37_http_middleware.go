//go:build ignore

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Middleware wraps a handler and returns a new one. Compare
// Rust's tower::Layer, which wraps one Service into
// another. http.Handler is the interface; Middleware needs
// no new type — a func value IS the contract, and tag
// below matches its signature exactly.
type Middleware func(http.Handler) http.Handler

// tag sets a response header, then delegates to the inner
// handler.
func tag(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Tag", "lesson37")
			next.ServeHTTP(w, r)
		})
}

// count closes over a shared counter, bumping it per
// request. Rust would keep the state in a Service struct
// field; here a closure captures it. A real server runs
// handlers concurrently, so guard this with sync/atomic or
// a mutex; this single-request demo does not.
func count(hits *int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				*hits++
				next.ServeHTTP(w, r)
			})
	}
}

func main() {
	base := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello")
		})

	// Chain outward: count wraps tag wraps base, so a request
	// enters count first, the same outermost-first order as a
	// tower Layer stack.
	hits := 0
	handler := count(&hits)(tag(base))

	// httptest records the response in memory, like lesson
	// 35's transport: no network, no timestamps, fully
	// deterministic.
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET", "http://lesson.test/", nil)
	handler.ServeHTTP(recorder, request)

	fmt.Printf("status: %d\n", recorder.Code)
	fmt.Printf("x-tag:  %s\n", recorder.Header().Get("X-Tag"))
	fmt.Printf("body:   %s\n", recorder.Body.String())
	fmt.Printf("hits:   %d\n", hits)
}

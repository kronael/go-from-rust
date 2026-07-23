//go:build ignore

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Middleware wraps a handler and returns a new one.
type Middleware func(http.Handler) http.Handler

// tag sets a response header, then calls next.
func tag(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Tag", "lesson37")
			next.ServeHTTP(w, r)
		})
}

// count bumps a shared counter per request.
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

	// Chain outward: count wraps tag wraps base.
	hits := 0
	handler := count(&hits)(tag(base))

	// httptest records the response in memory.
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET", "http://lesson.test/", nil)
	handler.ServeHTTP(recorder, request)

	fmt.Printf("status: %d\n", recorder.Code)
	fmt.Printf("x-tag:  %s\n", recorder.Header().Get("X-Tag"))
	fmt.Printf("body:   %s\n", recorder.Body.String())
	fmt.Printf("hits:   %d\n", hits)
}

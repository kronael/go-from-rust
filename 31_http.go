//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

type handlerTransport struct {
	handler http.Handler
}

func (transport handlerTransport) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	transport.handler.ServeHTTP(recorder, request)
	return recorder.Result(), nil
}

func main() {
	mux := http.NewServeMux()
	// Pattern captures {name}; PathValue reads it.
	mux.HandleFunc("GET /hello/{name}",
		func(
			response http.ResponseWriter,
			request *http.Request,
		) {
			fmt.Fprintf(response, "hello %s\n",
				request.PathValue("name"))
		})

	// In-memory transport: runnable without a real network.
	client := &http.Client{
		Transport: handlerTransport{handler: mux},
		Timeout:   time.Second,
	}
	response, err := client.Get("http://lesson.test/hello/Ana")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("response: %d %s", response.StatusCode, body)
}

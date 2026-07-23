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
	// The method/path pattern captures one segment as name. PathValue
	// retrieves it.
	mux.HandleFunc("GET /hello/{name}",
		func(response http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(response, "hello %s\n", request.PathValue("name"))
		})

	// A real service passes mux to http.Server. The in-memory transport keeps
	// this lesson runnable in the Playground while exercising the same handler.
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

	// Reuse clients and always close response bodies. In a service, configure
	// http.Server timeouts and call Shutdown for graceful HTTP server shutdown.
}

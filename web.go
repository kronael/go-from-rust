//go:build web

package main

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"go-from-rust/internal/webapp"
)

//go:embed README.md
//go:embed [0-9][0-9]_*.go
//go:embed web
var assets embed.FS

func main() {
	readme, err := assets.ReadFile("README.md")
	if err != nil {
		log.Fatalf("web: read README.md: %v", err)
	}
	lessons, err := webapp.LoadLessons(assets, readme)
	if err != nil {
		log.Fatalf("web: load lessons: %v", err)
	}
	webRoot, err := fs.Sub(assets, "web")
	if err != nil {
		log.Fatalf("web: sub filesystem for web assets: %v", err)
	}

	srv := webapp.NewServer(lessons, webRoot, webapp.Config{
		PlaygroundURL: os.Getenv("PLAYGROUND_URL"),
		UserAgent:     os.Getenv("PLAYGROUND_USER_AGENT"),
	})

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:3999"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("web: listen on %s: %v", addr, err)
	}
	server := &http.Server{
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("web: listening on %s (%d lessons loaded)", listener.Addr(), len(lessons))
	log.Fatal(server.Serve(listener))
}

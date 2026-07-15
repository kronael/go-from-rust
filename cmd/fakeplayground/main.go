// Command fakeplayground is a minimal stand-in for the Go Playground compile
// API, used by tests and scripts/smoke-web.sh. It never contacts the real
// Playground and only binds to loopback addresses.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type compileEvent struct {
	Message string `json:"Message"`
	Kind    string `json:"Kind"`
	Delay   int64  `json:"Delay"`
}

type compileResponse struct {
	Errors string         `json:"Errors"`
	Events []compileEvent `json:"Events"`
}

func handleCompile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	source := r.FormValue("body")

	dir, err := os.MkdirTemp("", "fakeplayground-*")
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	mainFile := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainFile, []byte(source), 0o600); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", mainFile)
	cmd.Env = append(os.Environ(), "GOCACHE="+gocacheDir())
	stdout, err := cmd.Output()

	resp := compileResponse{Events: []compileEvent{}}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			resp.Errors = string(exitErr.Stderr)
		} else {
			resp.Errors = err.Error()
		}
	} else if len(stdout) > 0 {
		resp.Events = append(resp.Events, compileEvent{Message: string(stdout), Kind: "stdout"})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func gocacheDir() string {
	if dir := os.Getenv("FAKEPLAYGROUND_GOCACHE"); dir != "" {
		return dir
	}
	return filepath.Join(os.TempDir(), "fakeplayground-gocache")
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/compile", handleCompile)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("fakeplayground: listen on %s: %v", addr, err)
	}
	server := &http.Server{Handler: mux}
	log.Printf("fakeplayground: listening on %s", listener.Addr())
	log.Fatal(server.Serve(listener))
}

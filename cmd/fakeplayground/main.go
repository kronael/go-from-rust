// Command fakeplayground is a minimal stand-in for the Go Playground compile
// API, used by tests and scripts/playtest.sh. It never contacts the real
// Playground and only binds to loopback addresses.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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
	if r.Form.Get("version") != "2" || r.Form.Get("withVet") != "true" {
		http.Error(w, "bad playground request", http.StatusBadRequest)
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

	// `go run` spawns compile and link children. Killing it alone
	// leaves them holding the output pipes, and Wait blocks on those
	// pipes long past the deadline. Give the run its own process
	// group, signal the group, and bound the wait that follows.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	cmd.WaitDelay = 2 * time.Second

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	resp := compileResponse{Events: []compileEvent{}}
	if err != nil {
		resp.Errors = stderr.String()
		if resp.Errors == "" {
			resp.Errors = err.Error()
		}
	} else {
		if stdout.Len() > 0 {
			resp.Events = append(resp.Events, compileEvent{Message: stdout.String(), Kind: "stdout"})
		}
		if stderr.Len() > 0 {
			resp.Events = append(resp.Events, compileEvent{Message: stderr.String(), Kind: "stderr"})
		}
	}

	// The public Playground returns JSON with a text/plain content type.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_ = json.NewEncoder(w).Encode(resp)
}

func gocacheDir() string {
	if dir := os.Getenv("FAKEPLAYGROUND_GOCACHE"); dir != "" {
		return dir
	}
	return filepath.Join(os.TempDir(), "fakeplayground-gocache")
}

func listenLoopback(addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok || !tcpAddr.IP.IsLoopback() {
		resolvedAddr := listener.Addr().String()
		_ = listener.Close()
		return nil, fmt.Errorf("resolved address %s is not loopback", resolvedAddr)
	}
	return listener, nil
}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/compile", handleCompile)

	listener, err := listenLoopback(addr)
	if err != nil {
		log.Fatalf("fakeplayground: listen on %s: %v", addr, err)
	}
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Printf("fakeplayground: listening on %s", listener.Addr())
	log.Fatal(server.Serve(listener))
}

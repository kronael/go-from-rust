//go:build ignore

package main

import (
	"log/slog"
	"os"
)

func main() {
	// This text handler writes to stdout, includes DEBUG, and removes time so the
	// lesson output is deterministic.
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return attr
		},
	})
	logger := slog.New(handler)

	// tracing fields map to stable key/value arguments, not assembled strings.
	logger.Debug("request received", "method", "GET", "attempt", 1)
	logger.Info("request complete",
		slog.String("method", "GET"), slog.Int("status", 200))

	// With returns a new logger carrying shared attributes; it is not a
	// tracing span.
	requestLogger := logger.With(slog.String("request_id", "req-42"))
	requestLogger.Warn("slow response", slog.Int("status", 503))
}

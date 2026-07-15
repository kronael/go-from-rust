//go:build ignore

package main

import (
	"log/slog"
	"os"
)

func main() {
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
	logger.Info("request complete", slog.String("method", "GET"), slog.Int("status", 200))

	// tracing span context maps to a logger carrying shared fields.
	requestLogger := logger.With(slog.String("request_id", "req-42"))
	requestLogger.Warn("slow response", slog.Int("status", 503))
}

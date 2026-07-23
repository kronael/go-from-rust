//go:build ignore

package main

import (
	"log/slog"
	"os"
)

func main() {
	// Text handler to stdout; drop time for determinism.
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(
				groups []string, attr slog.Attr,
			) slog.Attr {
				if attr.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return attr
			},
		})
	logger := slog.New(handler)

	// Log key/value fields, not assembled strings.
	logger.Debug("request received",
		"method", "GET", "attempt", 1)
	logger.Info("request complete",
		slog.String("method", "GET"), slog.Int("status", 200))

	// With derives a logger carrying shared attributes.
	requestLogger := logger.With(
		slog.String("request_id", "req-42"))
	requestLogger.Warn("slow response",
		slog.Int("status", 503))
}

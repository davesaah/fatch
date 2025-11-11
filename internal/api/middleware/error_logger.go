package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"gitlab.com/davesaah/fatch/internal/types"
)

func init() {
	// Set slog to use JSON output
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

// HandlerFuncWithErr is a handler that returns an error
type HandlerFuncWithErr func(w http.ResponseWriter, r *http.Request) *types.ErrorDetails

// Handler wraps HandlerFuncWithErr into an http.HandlerFunc
// LoggerMiddleware logs errors and request details
func Handler(h HandlerFuncWithErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if err := h(w, r); err != nil {
			slog.Error(err.Message,
				"trace", err.Trace,
				"method", r.Method,
				"url", r.URL.String(),
				"remote_addr", r.RemoteAddr,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		}
	}
}

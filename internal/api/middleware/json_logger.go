package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// JSONLogger is a Chi middleware that logs each request in JSON
func JSONLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Use a ResponseWriter wrapper to capture status code
		ww := &responseWriterWrapper{ResponseWriter: w, status: 200}

		next.ServeHTTP(ww, r)

		slog.Info("HTTP request",
			"method", r.Method,
			"url", r.URL.String(),
			"status", ww.status,
			"remote_addr", r.RemoteAddr,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// responseWriterWrapper wraps http.ResponseWriter to capture the status code
type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

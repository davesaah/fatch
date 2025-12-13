package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"gitlab.com/davesaah/fatch/internal/database"
	"gitlab.com/davesaah/fatch/internal/http/handlers"
)

// const (
// 	maxMemoryBytes  = 5 * 1024 * 1024
// 	avgLogSizeBytes = 200
// )
//
// var (
// 	maxBatchSize = maxMemoryBytes / avgLogSizeBytes
// 	logChan      = make(chan database.Log, 10000) // 10,000 logs peak
// )

func LoggerMiddleware(h *handlers.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use ResponseWriterWrapper to capture status code
			ww := &responseWriterWrapper{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(ww, r)
			url := r.URL.String()

			go func() {
				ctx := context.WithoutCancel(r.Context())

				err := h.Service.Log(ctx, &database.Log{
					Level:     "INFO",
					Service:   strings.Split(url, "/")[1], // extract service info
					Timestamp: time.Now(),
					LogData: map[string]any{
						"method":      r.Method,
						"url":         url,
						"status":      ww.status,
						"remote_addr": r.RemoteAddr,
						"userID":      ctx.Value("userID"),
					},
				})

				if err != nil {
					slog.Error(err.Error(), "origin", "LoggerMiddleware")
				}
			}()
		})
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture the status code
type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

// WriteHeader extracts the status code from responseWriterWrapper and sets it
// WriteHeader conforms to the interface of http.ResponseWriter
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

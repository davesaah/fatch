package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/davesaah/fatch/internal/database"
	"github.com/davesaah/fatch/internal/http/handlers"
	"github.com/davesaah/fatch/pubsub"
	"github.com/davesaah/fatch/types"
)

// HandlerFuncWithErr is a handler that returns an error
type HandlerFuncWithErr func(w http.ResponseWriter, r *http.Request) *types.ErrorDetails

// MakeHandler wraps HandlerFuncWithErr into an http.HandlerFunc and logs errors
func MakeHandler(hf HandlerFuncWithErr, h *handlers.Handler, ps *pubsub.PubSub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use ResponseWriterWrapper to capture status code
		ww := &responseWriterWrapper{ResponseWriter: w, status: http.StatusOK}

		url := r.URL.String()

		start := time.Now()
		err := hf(ww, r)

		logItem := database.Log{
			Level:     "INFO",
			Service:   strings.Split(url, "/")[1], // extract service info
			Timestamp: time.Now(),
			LogData: map[string]any{
				"method":      r.Method,
				"url":         url,
				"status":      ww.status,
				"remote_addr": r.RemoteAddr,
				"userID":      r.Context().Value("userID"),
				"duration":    time.Since(start).String(),
			},
		}

		if err != nil {
			logItem.Level = err.Level
			logItem.LogData["msg"] = err.Message
			logItem.LogData["trace"] = err.Trace
			logItem.LogData["duration"] = time.Since(start).String()
		}

		ps.Publish("logs", logItem)
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

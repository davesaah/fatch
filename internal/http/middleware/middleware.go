package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/davesaah/fatch/internal/database"
	"gitlab.com/davesaah/fatch/internal/http/handlers"
	"gitlab.com/davesaah/fatch/types"
)

func init() {
	// Set slog to use JSON output
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

// HandlerFuncWithErr is a handler that returns an error
type HandlerFuncWithErr func(w http.ResponseWriter, r *http.Request) *types.ErrorDetails

// MakeHandler wraps HandlerFuncWithErr into an http.HandlerFunc and logs errors
func MakeHandler(hf HandlerFuncWithErr, h *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()

		if err := hf(w, r); err != nil {
			go func() {
				ctx := context.WithoutCancel(r.Context())

				err := h.Service.Log(ctx, &database.Log{
					Level:     err.Level,
					Service:   strings.Split(url, "/")[1], // extract service info
					Timestamp: time.Now(),
					LogData: map[string]any{
						"msg":         err.Message,
						"trace":       err.Trace,
						"method":      r.Method,
						"url":         url,
						"remote_addr": r.RemoteAddr,
						"userID":      ctx.Value("userID"),
					},
				})

				if err != nil {
					slog.Error(err.Error(), "origin", "MakeHandler")
				}
			}()
		}
	}
}

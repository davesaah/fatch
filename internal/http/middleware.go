package internalHTTP

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/davesaah/fatch/internal/config"
	"gitlab.com/davesaah/fatch/types"
)

func init() {
	// Set slog to use JSON output
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

// HandlerFuncWithErr is a handler that returns an error
type HandlerFuncWithErr func(w http.ResponseWriter, r *http.Request) *types.ErrorDetails

// MakeHandler wraps HandlerFuncWithErr into an http.HandlerFunc and logs errors
func MakeHandler(h HandlerFuncWithErr) http.HandlerFunc {
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

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "No unauthorised access", http.StatusUnauthorized)
			slog.Default().Error(err.Error())
			return
		}

		// get jwtSecret from config
		jwtSecret, err := config.LoadJWTConfig()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Panic(err)
		}

		// parse into custom claims
		claims := &types.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims,
			func(t *jwt.Token) (any, error) {
				return jwtSecret, nil
			})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Attach userID from claims
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JSONLoggerMiddleware logs each request in JSON
func JSONLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Use ResponseWriterWrapper to capture status code
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

// WriteHeader extracts the status code from responseWriterWrapper and sets it
// WriteHeader conforms to the interface of http.ResponseWriter
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

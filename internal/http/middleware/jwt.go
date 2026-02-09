package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/davesaah/fatch/types"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				http.Error(w, "No unauthorised access", http.StatusUnauthorized)
				return
			}

			_jwtSecret := []byte(jwtSecret)
			if len(_jwtSecret) == 0 {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				slog.Error("JWT secret key not set", "origin", "JWTAuthMiddleware")
				return
			}

			// parse into custom claims
			claims := &types.Claims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims,
				func(t *jwt.Token) (any, error) {
					return _jwtSecret, nil
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
}

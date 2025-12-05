// Package middleware defines the custom middlewares for fatch api
package middleware

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/davesaah/fatch/config"
	"gitlab.com/davesaah/fatch/types"
)

// JWTAuth middleware validates JWT tokens
func JWTAuth(next http.Handler) http.Handler {
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
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (any, error) {
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

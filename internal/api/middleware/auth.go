package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/davidreturns08/fatch/internal/config"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth middleware validates JWT tokens
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// get jwtSecret from config
		jwtSecret, err := config.LoadJWTConfig()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Panic(err)
		}

		// parse into custom claims
		claims := &types.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
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

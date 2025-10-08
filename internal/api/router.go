package api

import (
	"net/http"

	"github.com/davidreturns08/fatch/internal/api/handlers"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", handlers.CreateUser)
	mux.HandleFunc("PATCH /users/passwd", handlers.ChangePassword)

	mux.HandleFunc("POST /auth/verify", handlers.VerifyPassword)

	return mux
}

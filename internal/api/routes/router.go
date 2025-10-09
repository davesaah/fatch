package routes

import (
	"github.com/davidreturns08/fatch/internal/api/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

// setupV1Routes sets up the routes for the v1 API.
func SetupV1Routes() *chi.Mux {
	r := chi.NewRouter()

	// setup middlewares
	r.Use(middleware.Logger)

	// API ROUTES
	r.Route("/v1", func(r chi.Router) {
		// Swagger documentation
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		r.Get("/swagger/doc.json", handlers.ServeDocFile)

		// Health check endpoint
		r.Get("/health", handlers.HealthCheck)

		// USER ROUTES
		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlers.CreateUser)
			r.Patch("/passwd", handlers.ChangePassword)
		})

		// AUTH ROUTES
		r.Route("/auth", func(r chi.Router) {
			r.Post("/verify", handlers.VerifyPassword)
		})
	})

	return r
}

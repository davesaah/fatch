package routes

import (
	"github.com/davidreturns08/fatch/internal/api/handlers"
	"github.com/davidreturns08/fatch/internal/api/middleware"
	"github.com/go-chi/chi/v5"

	httpSwagger "github.com/swaggo/http-swagger"
)

// setupV1Routes sets up the routes for the v1 API.
func SetupV1Routes() *chi.Mux {
	r := chi.NewRouter()

	// setup middlewares
	r.Use(middleware.JSONLogger)

	// API ROUTES
	r.Route("/v1", func(r chi.Router) {
		// Swagger documentation
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		r.Get("/swagger/doc.json", handlers.ServeDocFile)

		// Health check endpoint
		r.Get("/health", middleware.Handler(handlers.HealthCheck))

		// USER ROUTES
		r.Route("/users", func(r chi.Router) {
			r.Post("/", middleware.Handler(handlers.CreateUser))
			r.Patch("/passwd", middleware.Handler(handlers.ChangePassword))
		})

		// AUTH ROUTES
		r.Route("/auth", func(r chi.Router) {
			r.Post("/verify", middleware.Handler(handlers.VerifyPassword))
		})
	})

	return r
}

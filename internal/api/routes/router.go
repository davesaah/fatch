package routes

import (
	"os"
	"time"

	"github.com/davidreturns08/fatch/internal/api/handlers"
	"github.com/davidreturns08/fatch/internal/api/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupV1Routes sets up the routes for the v1 API.
func SetupV1Routes() *chi.Mux {
	r := chi.NewRouter()

	// setup middlewares
	r.Use(chiMiddleware.AllowContentType("application/json"))

	// update later
	r.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.JSONLogger)

	// limit to 100 requests per minute for each unique IP
	// look into expanding the rate limiter function to sensitive endpoints
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// add timeout to request
	r.Use(chiMiddleware.Timeout(time.Second * 1))

	// API ROUTES
	r.Route("/v1", func(r chi.Router) {
		if os.Getenv("ENVIRONMENT") == "dev" {
			// profiler
			r.Mount("/debug", chiMiddleware.Profiler())

			// Swagger documentation
			r.Get("/swagger/*", httpSwagger.WrapHandler)
			r.Get("/swagger/doc.json", handlers.ServeDocFile)
		}

		// Health check endpoint
		r.Get("/health", middleware.Handler(handlers.HealthCheck))

		// USER ROUTES
		r.Route("/users", func(r chi.Router) {
			r.Post("/", middleware.Handler(handlers.CreateUser))
		})

		// AUTH ROUTES
		r.Post("/auth/verify", middleware.Handler(handlers.VerifyPassword))

		// PROTECTED ROUTES: Requires authentication
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth)

			r.Patch("/auth/passwd", middleware.Handler(handlers.ChangePassword))
		})
	})

	return r
}

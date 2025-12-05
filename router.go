package main

import (
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"gitlab.com/davesaah/fatch/handlers"
	"gitlab.com/davesaah/fatch/middleware"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// setupV1Routes sets up the routes for the v1 API.
func setupV1Routes() *chi.Mux {
	origins := []string{"http://localhost:8000"}
	if os.Getenv("ENVIRONMENT") == "dev" {
		origins = append(origins, []string{"https://restfox.dev", "https://hoppscotch.io"}...)
	}

	r := chi.NewRouter()

	// setup middlewares
	r.Use(chiMiddleware.AllowContentType("application/json"))

	// update later
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
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
			r.Mount("/debug", chiMiddleware.Profiler()) // profiler

			// Swagger documentation
			r.Get("/swagger/*", httpSwagger.WrapHandler)
			r.Get("/swagger/doc.json", handlers.ServeDocFile)
		}

		r.Get("/health", middleware.Handler(handlers.HealthCheck))

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", middleware.Handler(handlers.VerifyPassword))
			r.Post("/register", middleware.Handler(handlers.CreateUser))
			// r.Post("/verify", middleware.Handler(handlers.VerifyUser))
		})

		// PROTECTED ROUTES: Requires authentication
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth)

			r.Patch("/users/passwd", middleware.Handler(handlers.ChangePassword))

			r.Route("/currencies", func(r chi.Router) {
				r.Get("/", middleware.Handler(handlers.GetAllCurrencies))
				r.Get("/{id}", middleware.Handler(handlers.GetCurrencyByID))
			})

			r.Route("/accounts", func(r chi.Router) {
				r.Post("/", middleware.Handler(handlers.CreateAccount))
				r.Get("/{id}", middleware.Handler(handlers.GetAccountByID))
				r.Get("/", middleware.Handler(handlers.GetAllUserAccounts))
				r.Patch("/{id}", middleware.Handler(handlers.ArchiveAccountByID))
			})
		})
	})

	return r
}

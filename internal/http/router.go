package internalHTTP

import (
	"net/http"
	"time"

	"github.com/davesaah/fatch/internal/http/handlers"
	"github.com/davesaah/fatch/internal/http/middleware"
	"github.com/davesaah/fatch/pubsub"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handlers.Handler, ps *pubsub.PubSub) http.Handler {
	origins := []string{"http://localhost:8000"}

	r := chi.NewRouter()

	// setup middlewares
	r.Use(chiMiddleware.AllowContentType("application/json"))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(chiMiddleware.Recoverer)

	// limit to 100 requests per minute for each unique IP
	// look into expanding the rate limiter function to sensitive endpoints
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// add timeout to request
	r.Use(chiMiddleware.Timeout(time.Second * 5))

	// API ROUTES
	if h.Config.Environment == "development" {
		r.Mount("/debug", chiMiddleware.Profiler()) // profiler

		// Swagger documentation
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		r.Get("/swagger/doc.json", h.ServeDocFile)
	}

	r.Get("/health", middleware.MakeHandler(h.HealthCheck, h, ps))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", middleware.MakeHandler(h.Login, h, ps))
		r.Post("/register", middleware.MakeHandler(h.Register, h, ps))
		r.Post("/verify", middleware.MakeHandler(h.VerifyUser, h, ps))

		// Protected auth routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuthMiddleware(h.Config.JWTSecret))
			r.Patch("/passwd", middleware.MakeHandler(h.ChangePassword, h, ps))
			r.Post("/logout", middleware.MakeHandler(h.Logout, h, ps))
			r.Delete("/delete", middleware.MakeHandler(h.DeleteUser, h, ps))
		})
	})

	// PROTECTED ROUTES: Requires authentication
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(h.Config.JWTSecret))

		r.Get("/currencies", middleware.MakeHandler(h.GetCurrencies, h, ps))

		r.Route("/categories", func(r chi.Router) {
			r.Post("/", middleware.MakeHandler(h.CreateCategory, h, ps))
			r.Get("/", middleware.MakeHandler(h.GetCategories, h, ps))
			r.Delete("/", middleware.MakeHandler(h.DeleteCategory, h, ps))
			r.Patch("/", middleware.MakeHandler(h.UpdateCategory, h, ps))
		})

		// r.Route("/subcategories", func(r chi.Router) {
		// 	r.Post("/", middleware.MakeHandler(h.CreateSubCategory, h, ps))
		// 	r.Get("/", middleware.MakeHandler(h.GetSubCategories, h, ps))
		// })

		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", middleware.MakeHandler(h.CreateAccount, h, ps))
			r.Get("/{id}", middleware.MakeHandler(h.GetAccountByID, h, ps))
			r.Get("/", middleware.MakeHandler(h.GetAllUserAccounts, h, ps))
			r.Patch("/{id}", middleware.MakeHandler(h.ArchiveAccountByID, h, ps))
		})
	})

	return r
}

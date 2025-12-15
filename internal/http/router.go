package internalHTTP

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"gitlab.com/davesaah/fatch/internal/http/handlers"
	"gitlab.com/davesaah/fatch/internal/http/middleware"
	"gitlab.com/davesaah/fatch/pubsub"

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
	r.Use(chiMiddleware.Timeout(time.Second * 1))

	// API ROUTES
	if os.Getenv("ENVIRONMENT") == "dev" {
		r.Mount("/debug", chiMiddleware.Profiler()) // profiler

		// Swagger documentation
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		r.Get("/swagger/doc.json", h.ServeDocFile)
	}

	r.Get("/health", middleware.MakeHandler(h.HealthCheck, h, ps))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", middleware.MakeHandler(h.Login, h, ps))
		r.Post("/register", middleware.MakeHandler(h.CreateUser, h, ps))
		// r.Post("/verify", middleware.Handler(h.VerifyUser))
	})

	// PROTECTED ROUTES: Requires authentication
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)

		r.Patch("/users/passwd", middleware.MakeHandler(h.ChangePassword, h, ps))

		r.Route("/currencies", func(r chi.Router) {
			r.Get("/", middleware.MakeHandler(h.GetAllCurrencies, h, ps))
			r.Get("/{id}", middleware.MakeHandler(h.GetCurrencyByID, h, ps))
		})

		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", middleware.MakeHandler(h.CreateAccount, h, ps))
			r.Get("/{id}", middleware.MakeHandler(h.GetAccountByID, h, ps))
			r.Get("/", middleware.MakeHandler(h.GetAllUserAccounts, h, ps))
			r.Patch("/{id}", middleware.MakeHandler(h.ArchiveAccountByID, h, ps))
		})
	})

	return r
}

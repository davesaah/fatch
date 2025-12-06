package internalHTTP

import (
	"net/http"
	"os"
	// "time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"gitlab.com/davesaah/fatch/internal/http/handlers"

	// "github.com/go-chi/httprate"

	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handlers.Handler) http.Handler {
	origins := []string{"http://localhost:8000"}

	r := chi.NewRouter()

	// setup middlewares
	r.Use(middleware.AllowContentType("application/json"))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.Recoverer)
	r.Use(JSONLoggerMiddleware)

	// limit to 100 requests per minute for each unique IP
	// look into expanding the rate limiter function to sensitive endpoints
	// r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// add timeout to request
	// r.Use(middleware.Timeout(time.Second * 1))

	// API ROUTES
	if os.Getenv("ENVIRONMENT") == "dev" {
		r.Mount("/debug", middleware.Profiler()) // profiler

		// Swagger documentation
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		r.Get("/swagger/doc.json", h.ServeDocFile)
	}

	r.Get("/health", MakeHandler(h.HealthCheck))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", MakeHandler(h.Login))
		r.Post("/register", MakeHandler(h.CreateUser))
		// r.Post("/verify", middleware.Handler(h.VerifyUser))
	})

	// PROTECTED ROUTES: Requires authentication
	r.Group(func(r chi.Router) {
		r.Use(JWTAuthMiddleware)

		r.Patch("/users/passwd", MakeHandler(h.ChangePassword))

		r.Route("/currencies", func(r chi.Router) {
			r.Get("/", MakeHandler(h.GetAllCurrencies))
			r.Get("/{id}", MakeHandler(h.GetCurrencyByID))
		})

		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", MakeHandler(h.CreateAccount))
			r.Get("/{id}", MakeHandler(h.GetAccountByID))
			r.Get("/", MakeHandler(h.GetAllUserAccounts))
			r.Patch("/{id}", MakeHandler(h.ArchiveAccountByID))
		})
	})

	return r
}

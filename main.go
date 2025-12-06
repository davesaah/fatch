package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/davesaah/fatch/internal/database"
	internalHTTP "gitlab.com/davesaah/fatch/internal/http"
	"gitlab.com/davesaah/fatch/internal/http/handlers"
	"gitlab.com/davesaah/fatch/internal/services"
)

// @title Fatch API
// @version 1.0
// @description Track your money; fetch insights on spending; budget effectively.
// @host localhost:8000
// @BasePath /
func main() {
	ctx := context.Background()
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	service := services.NewService(pool)
	handler := handlers.NewHandler(service)
	router := internalHTTP.NewRouter(handler)

	log.Println("API server started on http://localhost:8000")
	if os.Getenv("ENVIRONMENT") == "dev" {
		log.Println("API docs available at http://localhost:8000/swagger/index.html")
		log.Println("API profiler available at http://localhost:8000/debug")
	}

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server error:", err)
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	server.Shutdown(ctxShutdown)
	log.Println("Server stopped gracefully")
}

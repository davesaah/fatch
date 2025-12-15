package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"gitlab.com/davesaah/fatch/internal/database"
	internalHTTP "gitlab.com/davesaah/fatch/internal/http"
	"gitlab.com/davesaah/fatch/internal/http/handlers"
	"gitlab.com/davesaah/fatch/internal/services"
	"gitlab.com/davesaah/fatch/pubsub"
)

const (
	maxMemoryBytes  = 5 * 1024 * 1024
	avgLogSizeBytes = 200
)

var (
	maxBatchSize = maxMemoryBytes / avgLogSizeBytes
)

// @title Fatch API
// @version 1.0
// @description Track your money; fetch insights on spending; budget effectively.
// @host localhost:8000
// @BasePath /
func main() {
	ctx := context.Background()

	// create database pool
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// create pub/sub for realtime streaming
	ps := pubsub.New()
	sub := ps.Subscribe("logs", maxBatchSize)

	// Initialise API layers
	service := services.NewService(pool)

	// stream logs to DB
	go func() {
		for msg := range sub {
			service.Log(ctx, &msg)
		}
	}()

	handler := handlers.NewHandler(service)
	router := internalHTTP.NewRouter(handler, ps)

	log.Println("API server started on http://localhost:8000")
	if os.Getenv("ENVIRONMENT") == "dev" {
		log.Println("API docs available at http://localhost:8000/swagger/index.html")
		log.Println("API profiler available at http://localhost:8000/debug")
	}

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}

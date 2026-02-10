package main

import (
	"context"
	"log"
	"net/http"

	"github.com/davesaah/fatch/internal/config"
	"github.com/davesaah/fatch/internal/database"
	internalHTTP "github.com/davesaah/fatch/internal/http"
	"github.com/davesaah/fatch/internal/http/handlers"
	"github.com/davesaah/fatch/internal/services"
	"github.com/davesaah/fatch/pubsub"
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

// @contact.name   David Saah
// @contact.url    https://davesaah.com
// @contact.email  dave@davesaah.com

// @license.name GNU General Public License v3.0
// @license.url https://choosealicense.com/licenses/gpl-3.0/

// @host localhost:8000
// @BasePath /
func main() {
	ctx := context.Background()

	// load config
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// create database pool
	pool, err := database.NewPool(ctx, &config.DBConfig)
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

	handler := handlers.NewHandler(service, config)
	router := internalHTTP.NewRouter(handler, ps)

	log.Println("API server started on http://localhost:8000")
	if config.Environment == "development" {
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

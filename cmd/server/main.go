package main

import (
	"log"
	"net/http"

	"github.com/davidreturns08/fatch/internal/api"
)

func main() {
	mux := api.SetupRoutes()

	server := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	log.Printf("Starting server at localhost%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

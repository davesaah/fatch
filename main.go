package main

import (
	"log"
	"net/http"
	"os"
)

// @title Fatch API
// @version 1.0
// @description Track your money; fetch insights on spending; budget effectively.
// @host localhost:8000
// @BasePath /
func main() {
	mux := setupRoutes()

	log.Println("API server started on http://localhost:8000")
	if os.Getenv("ENVIRONMENT") == "dev" {
		log.Println("API docs available at http://localhost:8000/swagger/index.html")
		log.Println("API profiler available at http://localhost:8000/debug")
	}

	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

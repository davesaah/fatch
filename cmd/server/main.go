package main

import (
	"log"
	"net/http"

	"github.com/davidreturns08/fatch/internal/api/routes"
)

// @title Fatch API
// @version 1.0
// @description Track your money; fetch insights on spending; budget effectively.
// @host localhost:5000
// @BasePath /api/v1
func main() {
	mux := routes.SetupV1Routes()
	log.Println("API server started on http://localhost:500/api/v1")
	log.Println("API docs available at http://localhost:5000/swagger/index.html")
	http.ListenAndServe(":5000", mux)
}

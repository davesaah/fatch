package main

import (
	"log"
	"net/http"

	"github.com/davidreturns08/fatch/internal/api/routes"
)

// @title Fatch API
// @version 1.0
// @description Track your money; fetch insights on spending; budget effectively.
// @host api.fatch.laelfamily.org:8443
// @BasePath /v1
func main() {
	mux := routes.SetupV1Routes()

	log.Println("API server started on https://api.fatch.laelfamily.org:8443/v1")
	log.Println("API docs available at https://api.fatch.laelfamily.org:8443/v1/swagger/index.html")

	err := http.ListenAndServeTLS(":8443", "certs/cert.pem", "certs/key.pem", mux)
	if err != nil {
		log.Fatal(err)
	}
}

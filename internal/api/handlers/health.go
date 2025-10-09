package handlers

import (
	"net/http"

	"github.com/davidreturns08/fatch/internal/config"
	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/types"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// check if loading configurations are successful
	_, err := config.LoadDBConfig()
	if err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return
	}

	// check if database connection is established
	ctx := r.Context()
	db, err := database.NewConnection(ctx)
	if err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return
	}

	if err := db.Ping(ctx); err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return
	}

	types.ReturnJSON(w, types.OKResponse("Service is up", nil))
}

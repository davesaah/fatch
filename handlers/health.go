package handlers

import (
	"net/http"

	"gitlab.com/davesaah/fatch/config"
	"gitlab.com/davesaah/fatch/database"
	"gitlab.com/davesaah/fatch/types"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	// check if loading configurations are successful
	_, err := config.LoadDBConfig()
	if err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return &types.ErrorDetails{
			Trace:   err,
			Message: "Failed to load database config",
		}
	}

	// check if database connection is established
	ctx := r.Context()
	db, err := database.NewConnection(ctx)
	if err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return &types.ErrorDetails{
			Trace:   err,
			Message: "Failed to establish database connection",
		}
	}

	if err := db.Ping(ctx); err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return &types.ErrorDetails{
			Trace:   err,
			Message: "Failed to ping database",
		}
	}

	// check if jwt config is loaded successfully
	_, err = config.LoadJWTConfig()
	if err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return &types.ErrorDetails{
			Trace:   err,
			Message: "Failed to load jwt config",
		}
	}

	return types.ReturnJSON(w, types.OKResponse("Service is up", nil))
}

package handlers

import (
	"context"
	"net/http"

	"gitlab.com/davesaah/fatch/internal/config"
	"gitlab.com/davesaah/fatch/types"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := context.Background()
	// check if loading configurations are successful
	_, err := config.LoadDBConfig()
	if err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return &types.ErrorDetails{
			Trace:   err,
			Message: "Failed to load database config",
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

	// check db ping
	if err = h.Service.DB.Ping(ctx); err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return &types.ErrorDetails{
			Trace:   err,
			Message: "Failed to ping DB",
		}
	}

	return types.ReturnJSON(w, types.OKResponse("Service is up", nil))
}

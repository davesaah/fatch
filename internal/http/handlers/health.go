package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"gitlab.com/davesaah/fatch/types"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := context.Background()

	var env = []string{
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
		"DB_SCHEMA",
		"JWT_SECRET",
	}

	var missingValues []string

	for _, variable := range env {
		val := os.Getenv(variable)
		if val == "" {
			missingValues = append(missingValues, variable)
		}
	}

	// check if all required env values are set
	if len(missingValues) != 0 {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return &types.ErrorDetails{
			Trace:   fmt.Errorf("Missing env values: %+v", missingValues),
			Message: "Not all env variables are set",
			Level:   "ERROR",
		}
	}

	// check db ping
	if err := h.Service.DB.Ping(ctx); err != nil {
		types.ReturnJSON(w, types.ServiceUnavailableErrorResponse())
		return &types.ErrorDetails{
			Trace:   err,
			Level:   "ERROR",
			Message: "Failed to ping DB",
		}
	}

	return types.ReturnJSON(w, types.OKResponse("Service is up", nil))
}

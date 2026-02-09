package handlers

import (
	"net/http"

	"github.com/davesaah/fatch/types"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	return types.ReturnJSON(w, types.OKResponse("Service is up", nil))
}

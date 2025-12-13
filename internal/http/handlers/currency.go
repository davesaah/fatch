package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gitlab.com/davesaah/fatch/types"
)

func (h *Handler) GetCurrencyByID(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid currency id"))
		return &types.ErrorDetails{
			Message: "Unable to convert id to integer",
			Level:   "ERROR",
			Trace:   err,
		}
	}

	currency, errResponse, err := h.Service.GetCurrencyByID(r.Context(), id)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to get currency",
			Level:   "DEBUG",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("", currency))
}

func (h *Handler) GetAllCurrencies(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	currencies, errResponse, err := h.Service.GetAllCurrencies(r.Context())
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to get all currencies",
			Level:   "DEBUG",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("", currencies))
}

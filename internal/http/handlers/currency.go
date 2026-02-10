package handlers

import (
	"net/http"
	"strconv"

	"github.com/davesaah/fatch/types"
)

// GetCurrency
// @Summary Get all currencies or currency by ID
// @Description Get all currencies or currency by ID
// @Tags currencies
// @Accept json
// @Produce json
// @Param id query int false "Currency ID"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /currencies [get]
func (h *Handler) GetCurrencyByID(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	idStr := r.URL.Query().Get("id")
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

func (h *Handler) GetCurrencies(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	if r.URL.Query().Get("id") != "" {
		return h.GetCurrencyByID(w, r)
	}

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

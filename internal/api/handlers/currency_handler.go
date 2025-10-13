package handlers

import (
	"net/http"
	"strconv"

	"github.com/davidreturns08/fatch/internal/services"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/go-chi/chi/v5"
)

var currencyService *services.CurrencyService

func GetCurrencyById(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid currency id"))
		return &types.ErrorDetails{
			Message: "Unable to convert id to integer",
			Trace:   err,
		}
	}

	currency, errResponse, err := currencyService.GetCurrencyByID(r.Context(), id)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to get currency",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("", currency))
}

func GetAllCurrencies(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	currencies, errResponse, err := currencyService.GetAllCurrencies(r.Context())
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Fetch all currencies",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("", currencies))
}

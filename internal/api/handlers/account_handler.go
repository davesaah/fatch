package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/services"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var accountService *services.AccountService

func CreateAccount(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := r.Context()

	params := database.CreateAccountParams{
		Balance:     0.0,
		Description: "",
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
		return &types.ErrorDetails{
			Message: "Unable to parse json",
			Trace:   err,
		}
	}

	if params.AccountName == "" {
		return types.ReturnJSON(w, types.BadRequestErrorResponse("Account name cannot be empty"))
	}

	_, errResponse, err := currencyService.GetCurrencyByID(ctx, params.CurrencyID)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to get currency",
			Trace:   err,
		}
	}

	// extract user id from context
	userID := ctx.Value("userID").(pgtype.UUID)
	params.UserID = userID

	account, errResponse, err := accountService.CreateAccount(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to create account",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.CreatedResponse("Account created successfully", account))
}

func GetAccountByID(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid account id"))
		return &types.ErrorDetails{
			Message: "Unable to convert id to integer",
			Trace:   err,
		}
	}

	var params database.GetAccountDetailsParams
	params.AccountID = id

	// extract user id from context
	ctx := r.Context()
	userID := ctx.Value("userID").(pgtype.UUID)
	params.UserID = userID

	account, errResponse, err := accountService.GetAccountDetails(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to get account",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("Account retrieved successfully", account))
}

func GetAllUserAccounts(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	// extract user id from context
	ctx := r.Context()
	userID := ctx.Value("userID").(pgtype.UUID)

	accounts, errResponse, err := accountService.GetAllUserAccounts(ctx, userID)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to get all user accounts",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("", accounts))
}

func ArchiveAccountByID(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid account id"))
		return &types.ErrorDetails{
			Message: "Unable to convert id to integer",
			Trace:   err,
		}
	}

	var params database.ArchiveAccountByIDParams
	params.AccountID = id

	// extract user id from context
	ctx := r.Context()
	userID := ctx.Value("userID").(pgtype.UUID)
	params.UserID = userID

	errResponse, err := accountService.ArchiveAccount(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to archive account",
			Trace:   err,
		}
	}

	return types.ReturnJSON(w, types.OKResponse("Account archived successfully", nil))
}

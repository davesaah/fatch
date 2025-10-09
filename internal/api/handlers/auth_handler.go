package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/services"
	"github.com/davidreturns08/fatch/internal/types"
)

var authService services.AuthService

// ChangePassword handles changing a user's password.
func ChangePassword(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	var ctx = r.Context()
	var params database.ChangePasswordParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
	}

	// VALIDATE INPUT
	// Missing fields validation
	if params.OldPasswd == "" || params.NewPasswd == "" {
		return types.ReturnJSON(w,
			types.BadRequestErrorResponse("No empty fields allowed"),
		)
	}

	// password length validation
	if len(params.OldPasswd) < 8 || len(params.NewPasswd) < 8 {
		return types.ReturnJSON(w,
			types.PreconditionFailedErrorResponse("Password must be at least 8 characters long"),
		)
	}

	// call service to change password
	errResponse, err := authService.ChangePassword(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Failed to change password",
			Trace:   err,
		}
	}

	// Send success response
	return types.ReturnJSON(w, types.OKResponse("Password changed successfully", nil))
}

// VerifyPassword handles verifying a user's password.
func VerifyPassword(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	var ctx = r.Context()
	var params database.VerifyPasswordParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
	}

	// VALIDATE INPUT
	// Missing fields validation: Can sign in with email or username
	if (params.Email == "" && params.Username == "") || params.Passwd == "" {
		return types.ReturnJSON(w,
			types.BadRequestErrorResponse("No empty fields allowed"),
		)
	}

	// call service to create user
	response, errResponse, err := authService.VerifyPassword(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Failed to verify password",
			Trace:   err,
		}
	}

	// get user info
	user, errResponse, err := userService.GetUserById(ctx, database.GetUserByIdParams{
		UserID: response.UserID,
	})
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Failed to get user info",
			Trace:   err,
		}
	}

	responseMsg := fmt.Sprintf("%s logged in successfully", user.Username)

	// return success response
	return types.ReturnJSON(w, types.OKResponse(responseMsg, response))
}

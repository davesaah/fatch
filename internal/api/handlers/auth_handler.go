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
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var params database.ChangePasswordParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		types.BadRequestErrorResponse("Invalid JSON data")
		return
	}

	// VALIDATE INPUT
	// Missing fields validation
	if params.OldPasswd == "" || params.NewPasswd == "" {
		types.BadRequestErrorResponse("No empty fields allowed")
		return
	}

	// password length validation
	if len(params.OldPasswd) < 8 || len(params.NewPasswd) < 8 {
		types.ReturnJSON(w,
			types.PreconditionFailedErrorResponse("Password must be at least 8 characters long"),
		)
		return
	}

	// call service to change password
	errResponse := authService.ChangePassword(ctx, params)
	if errResponse != nil {
		types.ReturnJSON(w, errResponse)
		return
	}

	// Send success response
	types.ReturnJSON(w, types.OKResponse("Password changed successfully", nil))
}

// VerifyPassword handles verifying a user's password.
func VerifyPassword(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var params database.VerifyPasswordParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		types.BadRequestErrorResponse("Invalid JSON data")
		return
	}

	// VALIDATE INPUT
	// Missing fields validation: Can sign in with email or username
	if (params.Email == "" && params.Username == "") || params.Passwd == "" {
		types.BadRequestErrorResponse("No empty fields allowed")
		return
	}

	// call service to create user
	response, errResponse := authService.VerifyPassword(ctx, params)
	if errResponse != nil {
		types.ReturnJSON(w, errResponse)
		return
	}

	// get user info
	user, errResponse := userService.GetUserById(ctx, database.GetUserByIdParams{
		UserID: response.UserID,
	})
	if errResponse != nil {
		types.ReturnJSON(w, errResponse)
		return
	}

	responseMsg := fmt.Sprintf("%s logged in successfully", user.Username)

	// return success response
	types.ReturnJSON(w, types.OKResponse(responseMsg, response))
}

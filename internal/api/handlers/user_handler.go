package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/services"
	"github.com/davidreturns08/fatch/internal/types"
)

var userService services.UserService

// CreateUser handles the creation of a new user.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	// get & validate json data from request body
	var params database.CreateUserParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		types.WriteJSONError(w, &types.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON data",
		})
		return
	}

	// VALIDATE INPUT
	// Missing fields validation
	if params.Email == "" || params.Username == "" || params.Passwd == "" {
		types.WriteJSONError(w, &types.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "No empty fields allowed",
		})
		return
	}

	// password length validation
	if len(params.Passwd) < 8 {
		types.WriteJSONError(w, &types.ErrorResponse{
			StatusCode: http.StatusPreconditionFailed,
			Message:    "Password must be at least 8 characters long",
		})
		return
	}

	// TODO: Validate email with OTP

	// call service to create user
	errResponse := userService.CreateUser(ctx, params)
	if errResponse != nil {
		types.WriteJSONError(w, errResponse)
		return
	}

	// return success response
	types.WriteJSONSuccess(w, "User created successfully", nil)
}

package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/davesaah/fatch/database"
	"gitlab.com/davesaah/fatch/services"
	"gitlab.com/davesaah/fatch/types"
)

var userService services.UserService

// @Summary Register a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body database.CreateUserParams true "Request body for registering a new user"
// @Success 201 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 412 {object} types.ErrorResponse
// @Failure 409 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := r.Context()

	// get & validate json data from request body
	var params database.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
		return &types.ErrorDetails{
			Message: "Unable to parse json",
			Trace:   err,
		}
	}

	// VALIDATE INPUT
	// Missing fields validation
	if params.Email == "" || params.Username == "" || params.Passwd == "" {
		return types.ReturnJSON(w, types.BadRequestErrorResponse("No empty fields allowed"))
	}

	// password length validation
	if len(params.Passwd) < 8 {
		return types.ReturnJSON(w,
			types.PreconditionFailedErrorResponse("Password must be at least 8 characters long"),
		)
	}

	// TODO: Validate email with OTP

	// call service to create user
	errResponse, err := userService.CreateUser(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to create a new user",
			Trace:   err,
		}
	}

	// return success response
	return types.ReturnJSON(w, types.CreatedResponse("User created successfully", nil))
}

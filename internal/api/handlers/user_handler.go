package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/services"
	"github.com/davidreturns08/fatch/internal/types"
)

var userService services.UserService

// @Summary Create a new user
// @Description Registers a new user with a username, email, and password.
// @Tags users
// @Accept  json
// @Produce  json
// @Param request body database.CreateUserParams true "User registration data"
// @Success 200 {object} types.SuccessResponse "User created successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid JSON data or empty fields"
// @Failure 412 {object} types.ErrorResponse "Password must be at least 8 characters long"
// @Failure 409 {object} types.ErrorResponse "Username or email already exists"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	var ctx = r.Context()

	// get & validate json data from request body
	var params database.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
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

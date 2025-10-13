package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/davidreturns08/fatch/internal/config"
	"github.com/davidreturns08/fatch/internal/database"
	"github.com/davidreturns08/fatch/internal/services"
	"github.com/davidreturns08/fatch/internal/types"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/golang-jwt/jwt/v5"
)

var authService services.AuthService

// @Summary Change password for a user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.ChangePasswordParams true "Request body for changing user password"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 412 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /auth/passwd [post]
func ChangePassword(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := r.Context()
	var params types.ChangePasswordParams

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

	// don't allow same passwords
	if params.OldPasswd == params.NewPasswd {
		return types.ReturnJSON(w, types.PreconditionFailedErrorResponse("Old and new password fields must be different"))
	}

	// extract user id from context
	userID := ctx.Value("userID").(pgtype.UUID)

	// call service to change password
	errResponse, err := authService.ChangePassword(ctx, database.ChangePasswordParams{
		UserID:    userID,
		OldPasswd: params.OldPasswd,
		NewPasswd: params.NewPasswd,
	})
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

// @Summary Verify user login attempt
// @Tags auth
// @Accept json
// @Produce json
// @Param request body database.VerifyPasswordParams true "Request body for verifying user login attempt"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /auth/verify [post]
func VerifyPassword(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := r.Context()
	var params database.VerifyPasswordParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
	}

	// VALIDATE INPUT
	// Missing fields validation
	if params.Username == "" || params.Passwd == "" {
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
	user, errResponse, err := userService.GetUserById(ctx, response.UserID)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Failed to get user info",
			Trace:   err,
		}
	}

	responseMsg := fmt.Sprintf("%s logged in successfully", user.Username)

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &types.Claims{
		UserID: response.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// get jwtSecret from config
	jwtSecret, err := config.LoadJWTConfig()
	if err != nil {
		types.ReturnJSON(w, types.InternalServerErrorResponse())
		return &types.ErrorDetails{
			Message: "Failed to load jwt config",
			Trace:   err,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		types.ReturnJSON(w, types.InternalServerErrorResponse())
		return &types.ErrorDetails{
			Message: "Failed to create jwt token",
			Trace:   err,
		}
	}

	// Send token in header
	w.Header().Set("Authorization", "Bearer "+tokenString)

	// return success response
	return types.ReturnJSON(w, types.OKResponse(responseMsg, nil))
}

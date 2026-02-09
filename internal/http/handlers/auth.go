package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/davesaah/fatch/internal/database"
	"github.com/davesaah/fatch/types"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/golang-jwt/jwt/v5"
)

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
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := r.Context()
	var params database.ChangePasswordParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
		return &types.ErrorDetails{
			Message: "Unable to parse json",
			Level:   "ERROR",
			Trace:   err,
		}
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
			types.PreconditionFailedErrorResponse(
				"Password must be at least 8 characters long",
			))
	}

	// don't allow same passwords
	if params.OldPasswd == params.NewPasswd {
		return types.ReturnJSON(w, types.PreconditionFailedErrorResponse(
			"Old and new password fields must be different",
		))
	}

	// extract user id from context
	userID := ctx.Value("userID").(pgtype.UUID)
	params.UserID = userID

	// call service to change password
	errResponse, err := h.Service.ChangePassword(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Failed to change password",
			Level:   "WARN",
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
// @Param request body database.LoginParams true "Request body for verifying user login attempt"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := r.Context()
	var params database.LoginParams

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
		return &types.ErrorDetails{
			Message: "Unable to parse json",
			Level:   "ERROR",
			Trace:   err,
		}
	}

	// VALIDATE INPUT
	// Missing fields validation
	if params.Username == "" || params.Passwd == "" {
		return types.ReturnJSON(w,
			types.BadRequestErrorResponse("No empty fields allowed"),
		)
	}

	// call service to create user
	userID, errResponse, err := h.Service.Login(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Failed to verify password",
			Level:   "WARN",
			Trace:   err,
		}
	}

	// get user info
	user, errResponse, err := h.Service.GetUserByID(ctx, *userID)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Failed to get user info",
			Level:   "DEBUG",
			Trace:   err,
		}
	}

	responseMsg := fmt.Sprintf("%s logged in successfully", user.Username)

	expirationTime := time.Now().Add(1 * time.Minute)
	claims := &types.Claims{
		UserID: *userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtSecret := []byte(h.Config.JWTSecret)
	if len(jwtSecret) == 0 {
		types.ReturnJSON(w, types.InternalServerErrorResponse())
		return &types.ErrorDetails{
			Message: "JWT secret key not set",
			Level:   "ERROR",
			Trace:   nil,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		types.ReturnJSON(w, types.InternalServerErrorResponse())
		return &types.ErrorDetails{
			Message: "Failed to create jwt token",
			Level:   "ERROR",
			Trace:   err,
		}
	}

	// Set token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true --> for https
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(15 * time.Minute),
	})

	// return success response
	return types.ReturnJSON(w, types.OKResponse(responseMsg, nil))
}

// @Summary Register a new user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body database.RegisterParams true "Request body for registering a new user"
// @Success 201 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 412 {object} types.ErrorResponse
// @Failure 409 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := r.Context()

	// get & validate json data from request body
	var params database.RegisterUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
		return &types.ErrorDetails{
			Message: "Unable to parse json",
			Level:   "ERROR",
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
	otp, errResponse, err := h.Service.CreateUser(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to create a new user",
			Level:   "DEBUG",
			Trace:   err,
		}
	}

	// return success response
	return types.ReturnJSON(w, types.CreatedResponse(
		"User created successfully. Please verify your email",
		map[string]int{"otp": otp},
	))
}

func (h *Handler) VerifyUser(w http.ResponseWriter, r *http.Request) *types.ErrorDetails {
	ctx := r.Context()

	// get & validate json data from request body
	var params database.VerifyUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		types.ReturnJSON(w, types.BadRequestErrorResponse("Invalid JSON data"))
		return &types.ErrorDetails{
			Message: "Unable to parse json",
			Level:   "ERROR",
			Trace:   err,
		}
	}

	// VALIDATE INPUT
	// Missing fields validation
	if params.Username == "" || params.OTP == 0 || params.Passwd == "" {
		return types.ReturnJSON(w, types.BadRequestErrorResponse("No empty fields allowed"))
	}

	// call service to verify user
	errResponse, err := h.Service.VerifyUser(ctx, params)
	if err != nil {
		types.ReturnJSON(w, errResponse)
		return &types.ErrorDetails{
			Message: "Unable to verify user",
			Level:   "DEBUG",
			Trace:   err,
		}
	}

	// return success response
	return types.ReturnJSON(w, types.OKResponse("User verified successfully", nil))
}

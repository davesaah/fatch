package types

import (
	"encoding/json"
	"net/http"
)

// Common HTTP Status Codes
// 200 OK
// 201 Created
// 400 Bad Request
// 401 Unauthorized
// 403 Forbidden
// 404 Not Found
// 409 Conflict
// 412 Precondition Failed
// 500 Internal Server Error
// 503 Service Unavailable

func OKResponse(message string, data any) *SuccessResponse {
	return &SuccessResponse{
		Message:    message,
		Data:       data,
		StatusCode: http.StatusOK,
	}
}

func CreatedResponse(message string, data any) *SuccessResponse {
	return &SuccessResponse{
		Message:    message,
		Data:       data,
		StatusCode: http.StatusCreated,
	}
}

func BadRequestErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func UnauthorizedErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func ForbiddenErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NotFoundErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func ConflictErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func PreconditionFailedErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		StatusCode: http.StatusPreconditionFailed,
	}
}

func InternalServerErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}
}

func ServiceUnavailableErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Message:    "Service is down",
		StatusCode: http.StatusServiceUnavailable,
	}
}

// ReturnJSON writes a JSON response with the given status code, error type, and message.
func ReturnJSON(w http.ResponseWriter, resp Response) *ErrorDetails {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.GetStatusCode())

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return &ErrorDetails{
			Message: "Failed to encode response",
			Trace:   err,
		}
	}

	return nil
}

package types

import (
	"encoding/json"
	"net/http"
)

func InternalServerErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}
}

func BadRequestErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// WriteJSONError writes a JSON error response with the given status code, error type, and message.
func WriteJSONError(w http.ResponseWriter, resp *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}

// WriteJSONSuccess writes a JSON success response with the given message and data.
func WriteJSONSuccess(w http.ResponseWriter, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := SuccessResponse{
		Message:    message,
		Data:       data,
		StatusCode: http.StatusOK,
	}

	json.NewEncoder(w).Encode(resp)
}

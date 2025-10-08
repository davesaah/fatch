package types

// SuccessResponse represents a standardized success response.
type SuccessResponse struct {
	Message    string `json:"message"`
	StatusCode int
	Data       any `json:"data"`
}

// ErrorResponse represents a structured error response.
type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int
}

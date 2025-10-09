package types

// Response interface defines the contract for all response types.
type Response interface {
	GetStatusCode() int
}

// SuccessResponse represents a standardized success response.
type SuccessResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Data       any    `json:"data"`
}

func (s *SuccessResponse) GetStatusCode() int {
	return s.StatusCode
}

// ErrorResponse represents a structured error response.
type ErrorResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *ErrorResponse) GetStatusCode() int {
	return e.StatusCode
}

// ErrorDetails represents detailed error information.
type ErrorDetails struct {
	Message string `json:"msg"`
	Trace   error  `json:"trace"`
}

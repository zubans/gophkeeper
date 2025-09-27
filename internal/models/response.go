// Package models contains data structures used throughout the GophKeeper application.
package models

// APIResponse represents a standard API response.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
}

// NewSuccessResponse creates a new successful API response.
func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse creates a new error API response.
func NewErrorResponse(message string, code int) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	}
}

package client

import "fmt"

// LicenseChainError represents an error from the LicenseChain API
type LicenseChainError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

func (e *LicenseChainError) Error() string {
	return e.Message
}

// Error types
var (
	ErrInvalidAPIKey      = &LicenseChainError{Type: "invalid_api_key", Message: "Invalid API key provided"}
	ErrInvalidURL         = &LicenseChainError{Type: "invalid_url", Message: "Invalid URL"}
	ErrInvalidResponse    = &LicenseChainError{Type: "invalid_response", Message: "Invalid response from server"}
	ErrNetworkError       = &LicenseChainError{Type: "network_error", Message: "Network error occurred"}
	ErrValidationError    = &LicenseChainError{Type: "validation_error", Message: "Validation error"}
	ErrAuthenticationError = &LicenseChainError{Type: "authentication_error", Message: "Authentication error"}
	ErrNotFoundError      = &LicenseChainError{Type: "not_found_error", Message: "Resource not found"}
	ErrRateLimitError     = &LicenseChainError{Type: "rate_limit_error", Message: "Rate limit exceeded"}
	ErrServerError        = &LicenseChainError{Type: "server_error", Message: "Server error"}
	ErrUnknownError       = &LicenseChainError{Type: "unknown_error", Message: "Unknown error occurred"}
)

// NewValidationError creates a new validation error
func NewValidationError(message string) *LicenseChainError {
	return &LicenseChainError{
		Type:    "validation_error",
		Message: fmt.Sprintf("Validation error: %s", message),
	}
}

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(message string) *LicenseChainError {
	return &LicenseChainError{
		Type:    "authentication_error",
		Message: fmt.Sprintf("Authentication error: %s", message),
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *LicenseChainError {
	return &LicenseChainError{
		Type:    "not_found_error",
		Message: fmt.Sprintf("Not found: %s", message),
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string) *LicenseChainError {
	return &LicenseChainError{
		Type:    "rate_limit_error",
		Message: fmt.Sprintf("Rate limit exceeded: %s", message),
	}
}

// NewServerError creates a new server error
func NewServerError(message string) *LicenseChainError {
	return &LicenseChainError{
		Type:    "server_error",
		Message: fmt.Sprintf("Server error: %s", message),
	}
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(statusCode int, message string) *LicenseChainError {
	return &LicenseChainError{
		Type:    "http_error",
		Message: fmt.Sprintf("HTTP error %d: %s", statusCode, message),
		Code:    statusCode,
	}
}

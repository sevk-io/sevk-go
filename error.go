package sevk

import "fmt"

// Error represents an API error
type Error struct {
	StatusCode int
	Message    string
	Code       string
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("sevk: %d %s (code: %s)", e.StatusCode, e.Message, e.Code)
	}
	return fmt.Sprintf("sevk: %d %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found error
func (e *Error) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error
func (e *Error) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsForbidden returns true if the error is a 403 Forbidden error
func (e *Error) IsForbidden() bool {
	return e.StatusCode == 403
}

// IsBadRequest returns true if the error is a 400 Bad Request error
func (e *Error) IsBadRequest() bool {
	return e.StatusCode == 400
}

// IsServerError returns true if the error is a 5xx server error
func (e *Error) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

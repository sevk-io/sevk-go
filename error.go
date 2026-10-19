package sevk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Error represents an API error
type Error struct {
	cause        error
	StatusCode   int
	Message      string
	Code         string
	Details      any
	ResponseBody any
	// Fields holds field-level validation errors (422), field name -> messages.
	Fields map[string][]string
	// RetryAfter is the Retry-After header in seconds on a 429, when present.
	RetryAfter *int
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.StatusCode == 0 {
		return e.Message
	}
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

// IsConflict returns true if the error is a 409 Conflict error
func (e *Error) IsConflict() bool {
	return e.StatusCode == 409
}

// IsValidationError returns true for a 422 (or a 400 carrying field errors)
func (e *Error) IsValidationError() bool {
	return e.StatusCode == 422 ||
		(e.StatusCode == 400 && (e.Code == "VALIDATION_ERROR" || len(e.Fields) > 0))
}

func (e *Error) Unwrap() error {
	return e.cause
}

func newRequestError(cause error) *Error {
	code := "NETWORK_ERROR"
	var networkError net.Error
	if errors.Is(cause, context.Canceled) {
		code = "REQUEST_CANCELLED"
	} else if errors.As(cause, &networkError) && networkError.Timeout() {
		code = "TIMEOUT"
	}
	return &Error{Message: cause.Error(), Code: code, cause: cause}
}

// IsRateLimited returns true if the error is a 429 Too Many Requests error
func (e *Error) IsRateLimited() bool {
	return e.StatusCode == 429
}

// IsServerError returns true if the error is a 5xx server error
func (e *Error) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// parseAPIError turns a non-2xx HTTP response into a typed *Error. It unwraps
// the nested `{ "error": { ... } }` envelope when present, normalizes the
// validation `fields` map to lists of strings, and reads the Retry-After header.
func parseAPIError(statusCode int, header http.Header, body []byte) *Error {
	apiErr := &Error{
		StatusCode: statusCode,
		Code:       "UNKNOWN_ERROR",
		RetryAfter: parseRetryAfter(header.Get("Retry-After")),
	}

	var decoded any
	if err := json.Unmarshal(body, &decoded); err != nil {
		apiErr.ResponseBody = string(body)
		apiErr.Message = strings.TrimSpace(string(body))
		if apiErr.Message == "" {
			apiErr.Message = fmt.Sprintf("HTTP %d", statusCode)
		}
		return apiErr
	}
	apiErr.ResponseBody = decoded

	payload, _ := decoded.(map[string]any)
	if nested, ok := payload["error"].(map[string]any); ok {
		payload = nested
	}

	if message, ok := payload["message"].(string); ok && message != "" {
		apiErr.Message = message
	} else if text, ok := decoded.(string); ok {
		apiErr.Message = strings.TrimSpace(text)
	}
	if apiErr.Message == "" {
		apiErr.Message = fmt.Sprintf("HTTP %d", apiErr.StatusCode)
	}
	if code, ok := payload["code"].(string); ok && code != "" {
		apiErr.Code = code
	}
	apiErr.Fields = normalizeFields(payload["fields"])
	apiErr.Details = payload["details"]
	return apiErr
}

// normalizeFields coerces each field value (a string or a list of strings on
// the wire) into a list of strings.
func normalizeFields(raw any) map[string][]string {
	fields, ok := raw.(map[string]any)
	if !ok || len(fields) == 0 {
		return nil
	}
	out := make(map[string][]string, len(fields))
	for key, value := range fields {
		switch typed := value.(type) {
		case []any:
			list := make([]string, len(typed))
			for i, item := range typed {
				list[i] = fmt.Sprint(item)
			}
			out[key] = list
		case string:
			out[key] = []string{typed}
		default:
			out[key] = []string{fmt.Sprint(typed)}
		}
	}
	return out
}

const maxRetryAfterSeconds = 2_147_483_647

// parseRetryAfter parses an integer-seconds Retry-After header. HTTP-date forms
// are ignored (returns nil), matching the other SDKs.
func parseRetryAfter(value string) *int {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return nil
		}
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return nil
	}
	seconds := int(parsed)
	if seconds > maxRetryAfterSeconds {
		return nil
	}
	return &seconds
}

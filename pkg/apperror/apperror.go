// Package apperror provides structured, typed application errors with
// optional stack traces for debug mode. These errors map cleanly to
// HTTP status codes and the dto.ErrorResponse format.
package apperror

import (
	"fmt"
	"net/http"
	"runtime"
)

// AppError is a structured error that carries an error code, user-facing
// message, optional developer details, and a captured stack trace.
type AppError struct {
	// Code is a short machine-readable error code (e.g. "NOT_FOUND")
	Code string
	// Message is a safe, user-facing error description
	Message string
	// Details contains extra info for developers (hidden from users in production)
	Details string
	// HTTPStatus is the recommended HTTP status code for this error
	HTTPStatus int
	// Err is the underlying wrapped error (if any)
	Err error
	// Stack is the captured stack trace at creation time
	Stack string
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap supports errors.Is / errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// captureStack returns the current goroutine stack trace as a string
func captureStack() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// --- Constructors ---

// NewNotFound creates a 404 Not Found error
func NewNotFound(resource string, err error) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    resource + " not found",
		HTTPStatus: http.StatusNotFound,
		Err:        err,
		Stack:      captureStack(),
	}
}

// NewValidation creates a 400 Bad Request error for validation failures
func NewValidation(message string, details string) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		Details:    details,
		HTTPStatus: http.StatusBadRequest,
		Stack:      captureStack(),
	}
}

// NewBadRequest creates a generic 400 Bad Request error
func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Err:        err,
		Stack:      captureStack(),
	}
}

// NewInternal creates a 500 Internal Server Error
func NewInternal(message string, err error) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
		Stack:      captureStack(),
	}
}

// NewDatabase creates a 500 error specifically for database failures
func NewDatabase(message string, err error) *AppError {
	return &AppError{
		Code:       "DATABASE_ERROR",
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
		Stack:      captureStack(),
	}
}

// NewTimeout creates a 408 Request Timeout error
func NewTimeout(message string) *AppError {
	return &AppError{
		Code:       "TIMEOUT",
		Message:    message,
		HTTPStatus: http.StatusRequestTimeout,
		Stack:      captureStack(),
	}
}

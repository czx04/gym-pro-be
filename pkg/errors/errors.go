package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeValidation          = "VALIDATION_ERROR"
	ErrCodeInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError       = "DATABASE_ERROR"
	ErrCodeInvalidCredentials  = "INVALID_CREDENTIALS"
	ErrCodeTokenExpired        = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid        = "TOKEN_INVALID"
)

// NewAppError creates a new AppError
func NewAppError(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) *AppError {
	return &AppError{
		Code:    ErrCodeBadRequest,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    ErrCodeUnauthorized,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string) *AppError {
	return &AppError{
		Code:    ErrCodeForbidden,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

// NotFound creates a 404 Not Found error
func NotFound(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Status:  http.StatusNotFound,
	}
}

// Conflict creates a 409 Conflict error
func Conflict(message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
		Status:  http.StatusConflict,
	}
}

// Validation creates a 422 Validation error
func Validation(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
		Status:  http.StatusUnprocessableEntity,
	}
}

// InternalServer creates a 500 Internal Server Error
func InternalServer(message string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeInternalServer,
		Message: message,
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// DatabaseError creates a database error
func DatabaseError(operation string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeDatabaseError,
		Message: fmt.Sprintf("database error during %s", operation),
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// InvalidCredentials creates an invalid credentials error
func InvalidCredentials() *AppError {
	return &AppError{
		Code:    ErrCodeInvalidCredentials,
		Message: "invalid email or password",
		Status:  http.StatusUnauthorized,
	}
}

// TokenExpired creates a token expired error
func TokenExpired() *AppError {
	return &AppError{
		Code:    ErrCodeTokenExpired,
		Message: "token has expired",
		Status:  http.StatusUnauthorized,
	}
}

// TokenInvalid creates a token invalid error
func TokenInvalid() *AppError {
	return &AppError{
		Code:    ErrCodeTokenInvalid,
		Message: "token is invalid",
		Status:  http.StatusUnauthorized,
	}
}

// Wrap wraps an error with a message
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}
	
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	
	return &AppError{
		Code:    ErrCodeInternalServer,
		Message: message,
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// Is checks if an error is a specific AppError
func Is(err error, code string) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

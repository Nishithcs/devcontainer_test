package errors

import (
	"fmt"
	"net/http"
	"strings"
)

// ErrorType represents different types of errors
type ErrorType string

const (
	ErrorTypeInternal   ErrorType = "INTERNAL"
	ErrorTypeValidation ErrorType = "VALIDATION"
	ErrorTypeNotFound   ErrorType = "NOT_FOUND"
	ErrorTypeAuth       ErrorType = "AUTHENTICATION"
	ErrorTypeForbidden  ErrorType = "FORBIDDEN"
	ErrorTypeBadRequest ErrorType = "BAD_REQUEST"
)

// AppError represents a structured application error
type AppError struct {
	Type     ErrorType           `json:"-"`                // Internal type for categorization
	Code     string              `json:"code"`             // Error code for clients
	Message  string              `json:"message"`          // User-friendly message
	Detail   string              `json:"detail"`           // Detailed error information
	HTTPCode int                 `json:"-"`                // HTTP status code
	Raw      error               `json:"-"`                // Original error
	Metadata any                 `json:"metadata"`         // Additional error context
	Errors   map[string][]string `json:"errors,omitempty"` // Validation errors
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}

// NewError creates a new AppError
func NewError(errType ErrorType, code string, msg string, raw error) *AppError {
	httpCode := getHTTPCode(errType)
	detail := msg
	if raw != nil {
		detail = raw.Error()
	}

	return &AppError{
		Type:     errType,
		Code:     code,
		Message:  msg,
		Detail:   detail,
		HTTPCode: httpCode,
		Raw:      raw,
	}
}

// Helper functions for common errors
func NewInternalError(code string, err error) *AppError {
	return NewError(ErrorTypeInternal, code, "Internal server error", err)
}

func NewValidationError(msg string, errors map[string][]string) *AppError {
	return &AppError{
		Type:     ErrorTypeValidation,
		Code:     "VALIDATION",
		Message:  msg,
		Detail:   msg,
		HTTPCode: getHTTPCode(ErrorTypeValidation),
		Errors:   errors,
	}
}

func NewNotFoundError(resource string) *AppError {
	return NewError(ErrorTypeNotFound, "NOT_FOUND",
		fmt.Sprintf("%s not found", strings.Title(resource)), nil)
}

func NewAuthenticationError(msg string) *AppError {
	return NewError(ErrorTypeAuth, "UNAUTHORIZED", msg, nil)
}

func NewForbiddenError(msg string) *AppError {
	return NewError(ErrorTypeForbidden, "FORBIDDEN", msg, nil)
}

// getHTTPCode maps error types to HTTP status codes
func getHTTPCode(errType ErrorType) int {
	fmt.Println(fmt.Sprintf("Mapping error type %s to HTTP code", errType))
	switch errType {
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	case ErrorTypeValidation:
		return http.StatusUnprocessableEntity
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeAuth:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

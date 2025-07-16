package errors

import "fmt"

// ErrorType represents the type of error.
type ErrorType int

const (
	// ErrorTypeUnknown represents an unknown error.
	ErrorTypeUnknown ErrorType = iota
	// ErrorTypeAuthentication represents authentication errors.
	ErrorTypeAuthentication
	// ErrorTypeRateLimit represents rate limit errors.
	ErrorTypeRateLimit
	// ErrorTypeNotFound represents not found errors.
	ErrorTypeNotFound
	// ErrorTypeParse represents parsing errors.
	ErrorTypeParse
	// ErrorTypeNetwork represents network errors.
	ErrorTypeNetwork
	// ErrorTypeValidation represents validation errors.
	ErrorTypeValidation
)

// TypedError is an error with a type.
type TypedError interface {
	error
	Type() ErrorType
}

// typedError implements TypedError.
type typedError struct {
	errorType ErrorType
	message   string
	cause     error
}

func (e *typedError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

func (e *typedError) Type() ErrorType {
	return e.errorType
}

func (e *typedError) Unwrap() error {
	return e.cause
}

// NewAuthenticationError creates a new authentication error.
func NewAuthenticationError(message string) error {
	return &typedError{
		errorType: ErrorTypeAuthentication,
		message:   message,
	}
}

// NewRateLimitError creates a new rate limit error.
func NewRateLimitError(message string) error {
	return &typedError{
		errorType: ErrorTypeRateLimit,
		message:   message,
	}
}

// NewNotFoundError creates a new not found error.
func NewNotFoundError(message string) error {
	return &typedError{
		errorType: ErrorTypeNotFound,
		message:   message,
	}
}

// NewParseError creates a new parse error.
func NewParseError(message string) error {
	return &typedError{
		errorType: ErrorTypeParse,
		message:   message,
	}
}

// NewNetworkError creates a new network error.
func NewNetworkError(message string) error {
	return &typedError{
		errorType: ErrorTypeNetwork,
		message:   message,
	}
}

// NewValidationError creates a new validation error.
func NewValidationError(message string) error {
	return &typedError{
		errorType: ErrorTypeValidation,
		message:   message,
	}
}

// WrapError wraps an error with a type and message.
func WrapError(err error, errorType ErrorType, message string) error {
	return &typedError{
		errorType: errorType,
		message:   message,
		cause:     err,
	}
}

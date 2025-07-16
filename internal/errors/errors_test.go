package errors_test

import (
	"errors"
	"testing"

	semvererrors "github.com/binary-install/semver/internal/errors"
)

func TestErrorType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantType semvererrors.ErrorType
	}{
		{
			name:     "authentication error",
			err:      semvererrors.NewAuthenticationError("invalid token"),
			wantType: semvererrors.ErrorTypeAuthentication,
		},
		{
			name:     "rate limit error",
			err:      semvererrors.NewRateLimitError("API rate limit exceeded"),
			wantType: semvererrors.ErrorTypeRateLimit,
		},
		{
			name:     "not found error",
			err:      semvererrors.NewNotFoundError("repository not found"),
			wantType: semvererrors.ErrorTypeNotFound,
		},
		{
			name:     "parse error",
			err:      semvererrors.NewParseError("invalid version format"),
			wantType: semvererrors.ErrorTypeParse,
		},
		{
			name:     "network error",
			err:      semvererrors.NewNetworkError("connection timeout"),
			wantType: semvererrors.ErrorTypeNetwork,
		},
		{
			name:     "validation error",
			err:      semvererrors.NewValidationError("invalid constraint"),
			wantType: semvererrors.ErrorTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var typedErr semvererrors.TypedError
			if !errors.As(tt.err, &typedErr) {
				t.Fatalf("error does not implement TypedError interface")
			}

			if got := typedErr.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}

			if msg := tt.err.Error(); msg == "" {
				t.Error("Error() returned empty string")
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := semvererrors.WrapError(baseErr, semvererrors.ErrorTypeNetwork, "network operation failed")

	if !errors.Is(wrappedErr, baseErr) {
		t.Error("wrapped error should contain base error")
	}

	var typedErr semvererrors.TypedError
	if !errors.As(wrappedErr, &typedErr) {
		t.Fatal("wrapped error should implement TypedError")
	}

	if typedErr.Type() != semvererrors.ErrorTypeNetwork {
		t.Errorf("Type() = %v, want %v", typedErr.Type(), semvererrors.ErrorTypeNetwork)
	}
}

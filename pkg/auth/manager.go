package auth

import (
	"os"
	"strings"

	semvererrors "github.com/binary-install/semver/internal/errors"
	"github.com/zalando/go-keyring"
)

const (
	// Environment variables for GitHub token
	envGitHubToken = "GITHUB_TOKEN"
	envGHToken     = "GH_TOKEN"

	// Keyring service name
	keyringService = "semver-resolver"
	keyringAccount = "github-token"
)

// manager implements TokenManager interface.
type manager struct{}

// NewTokenManager creates a new token manager.
func NewTokenManager() TokenManager {
	return &manager{}
}

// GetToken returns the GitHub token.
func (m *manager) GetToken() (string, error) {
	// Check environment variables first
	if token := os.Getenv(envGitHubToken); token != "" {
		return token, nil
	}
	if token := os.Getenv(envGHToken); token != "" {
		return token, nil
	}

	// Check system keyring
	token, err := keyring.Get(keyringService, keyringAccount)
	if err == nil && token != "" {
		return token, nil
	}

	// If keyring error is not "not found", wrap and return it
	if err != nil && !isKeyringNotFoundError(err) {
		return "", semvererrors.WrapError(err, semvererrors.ErrorTypeAuthentication, "failed to access keyring")
	}

	return "", semvererrors.NewAuthenticationError("no GitHub token found in environment variables or keyring")
}

// SetToken stores the token in the system keyring.
func (m *manager) SetToken(token string) error {
	if token == "" {
		return semvererrors.NewValidationError("token cannot be empty")
	}

	if err := keyring.Set(keyringService, keyringAccount, token); err != nil {
		return semvererrors.WrapError(err, semvererrors.ErrorTypeAuthentication, "failed to store token in keyring")
	}

	return nil
}

// DeleteToken removes the token from the system keyring.
func (m *manager) DeleteToken() error {
	err := keyring.Delete(keyringService, keyringAccount)
	if err != nil && !isKeyringNotFoundError(err) {
		return semvererrors.WrapError(err, semvererrors.ErrorTypeAuthentication, "failed to delete token from keyring")
	}
	return nil
}

// isKeyringNotFoundError checks if the error is a "not found" error from keyring.
func isKeyringNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "The specified item could not be found") ||
		strings.Contains(errStr, "secret not found")
}

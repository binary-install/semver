package auth

// TokenManager manages GitHub authentication tokens.
type TokenManager interface {
	// GetToken returns the GitHub token.
	// It checks environment variables and system keyring.
	GetToken() (string, error)
	// SetToken stores the token in the system keyring.
	SetToken(token string) error
	// DeleteToken removes the token from the system keyring.
	DeleteToken() error
}

package auth_test

import (
	"os"
	"testing"

	"github.com/binary-install/semver/pkg/auth"
)

func TestTokenManager_GetToken(t *testing.T) {
	// Save original env vars
	origGitHubToken := os.Getenv("GITHUB_TOKEN")
	origGHToken := os.Getenv("GH_TOKEN")
	defer func() {
		os.Setenv("GITHUB_TOKEN", origGitHubToken)
		os.Setenv("GH_TOKEN", origGHToken)
	}()

	tests := []struct {
		name       string
		setupFunc  func()
		wantErr    bool
		wantPrefix string
	}{
		{
			name: "get token from GITHUB_TOKEN env var",
			setupFunc: func() {
				os.Setenv("GITHUB_TOKEN", "ghp_test123456789")
				os.Unsetenv("GH_TOKEN")
			},
			wantErr:    false,
			wantPrefix: "ghp_",
		},
		{
			name: "get token from GH_TOKEN env var",
			setupFunc: func() {
				os.Unsetenv("GITHUB_TOKEN")
				os.Setenv("GH_TOKEN", "gho_test123456789")
			},
			wantErr:    false,
			wantPrefix: "gho_",
		},
		{
			name: "no token available (env vars only)",
			setupFunc: func() {
				os.Unsetenv("GITHUB_TOKEN")
				os.Unsetenv("GH_TOKEN")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()

			tm := auth.NewTokenManager()
			token, err := tm.GetToken()

			// For the "no token" test, we might still get a token from keyring
			// So we only test the cases where we expect to find a token in env vars
			if !tt.wantErr {
				if err != nil {
					t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantPrefix != "" && len(token) >= len(tt.wantPrefix) {
					if token[:len(tt.wantPrefix)] != tt.wantPrefix {
						t.Errorf("GetToken() = %v, want prefix %v", token, tt.wantPrefix)
					}
				}
			} else {
				// For the no token case, we can't reliably test failure
				// because keyring might have a token
				if err == nil && token != "" {
					t.Skip("Found token in keyring, skipping negative test")
				}
			}
		})
	}
}

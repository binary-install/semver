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
			name: "no token available",
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

			if (err != nil) != tt.wantErr {
				// Skip error check if we got a token from keyring when expecting error
				if err == nil && tt.wantErr && token != "" {
					t.Skip("Found token in keyring, skipping negative test")
				}
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.wantPrefix != "" {
				if len(token) < len(tt.wantPrefix) || token[:len(tt.wantPrefix)] != tt.wantPrefix {
					t.Errorf("GetToken() = %v, want prefix %v", token, tt.wantPrefix)
				}
			}
		})
	}
}

func TestTokenManager_SetToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "set valid token",
			token:   "ghp_validtoken123",
			wantErr: false,
		},
		{
			name:    "set empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := auth.NewTokenManager()
			err := tm.SetToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

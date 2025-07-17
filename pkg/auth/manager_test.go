// +build integration

package auth_test

import (
	"testing"

	"github.com/binary-install/semver/pkg/auth"
)

func TestTokenManager_SetToken_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tm := auth.NewTokenManager()

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
			err := tm.SetToken(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTokenManager_DeleteToken_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tm := auth.NewTokenManager()

	// First, set a test token
	testToken := "ghp_test_delete_token_12345"
	err := tm.SetToken(testToken)
	if err != nil {
		t.Fatalf("Failed to set test token: %v", err)
	}

	// Delete the token
	err = tm.DeleteToken()
	if err != nil {
		t.Fatalf("DeleteToken() error = %v", err)
	}

	// Try to delete again (should not error even if token doesn't exist)
	err = tm.DeleteToken()
	if err != nil {
		t.Errorf("DeleteToken() on non-existent token error = %v", err)
	}
}

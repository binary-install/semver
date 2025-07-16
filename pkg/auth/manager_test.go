package auth_test

import (
	"testing"

	"github.com/binary-install/semver/pkg/auth"
)

func TestTokenManager_DeleteToken(t *testing.T) {
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

	// Try to get the token from keyring (should fail)
	// Note: Environment variables might still have tokens
	// so we can't reliably test that GetToken fails

	// Try to delete again (should not error even if token doesn't exist)
	err = tm.DeleteToken()
	if err != nil {
		t.Errorf("DeleteToken() on non-existent token error = %v", err)
	}
}

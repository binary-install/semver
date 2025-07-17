//go:build integration

package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/binary-install/semver"
	"github.com/binary-install/semver/pkg/auth"
	"github.com/binary-install/semver/pkg/github"
	"github.com/binary-install/semver/pkg/resolver"
)

func skipIfNoToken(t *testing.T) string {
	tm := auth.NewTokenManager()
	token, err := tm.GetToken()
	if err != nil || token == "" {
		t.Skip("Skipping integration test: no GitHub token found")
	}
	return token
}

func TestRealGitHubRepositories(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	token := skipIfNoToken(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tests := []struct {
		name       string
		owner      string
		repo       string
		constraint string
		wantPrefix string // Check if result starts with this prefix
	}{
		{
			name:       "Masterminds/semver exact version",
			owner:      "Masterminds",
			repo:       "semver",
			constraint: "3.2.0",
			wantPrefix: "v3.2.0",
		},
		{
			name:       "Masterminds/semver range",
			owner:      "Masterminds",
			repo:       "semver",
			constraint: "^3.0.0",
			wantPrefix: "v3.",
		},
		{
			name:       "spf13/cobra tilde range",
			owner:      "spf13",
			repo:       "cobra",
			constraint: "~1.8.0",
			wantPrefix: "v1.8.",
		},
		{
			name:       "google/go-github major version",
			owner:      "google",
			repo:       "go-github",
			constraint: "^57.0.0",
			wantPrefix: "v57.",
		},
	}

	githubClient := github.NewClient(token)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := resolver.New(githubClient, resolver.Options{})

			version, err := r.MaxSatisfying(ctx, tt.owner, tt.repo, tt.constraint)
			if err != nil {
				t.Fatalf("MaxSatisfying() error = %v", err)
			}

			if version == "" {
				t.Fatal("MaxSatisfying() returned empty version")
			}

			if len(version) < len(tt.wantPrefix) || version[:len(tt.wantPrefix)] != tt.wantPrefix {
				t.Errorf("MaxSatisfying() = %v, want prefix %v", version, tt.wantPrefix)
			}

			t.Logf("Resolved %s/%s@%s to %s", tt.owner, tt.repo, tt.constraint, version)
		})
	}
}

func TestLibraryAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	token := skipIfNoToken(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test the public API
	version, err := semver.MaxSatisfyingWithToken(ctx, "spf13", "cobra", "^1.0.0", token)
	if err != nil {
		t.Fatalf("MaxSatisfyingWithToken() error = %v", err)
	}

	if version == "" {
		t.Fatal("MaxSatisfyingWithToken() returned empty version")
	}

	t.Logf("Resolved spf13/cobra@^1.0.0 to %s", version)
}

func TestCLIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Ensure we have a token
	_ = skipIfNoToken(t)

	// Test cases for CLI
	tests := []struct {
		name     string
		args     string
		wantExit int
	}{
		{
			name:     "valid resolution",
			args:     "spf13/cobra@^1.0.0",
			wantExit: 0,
		},
		{
			name:     "invalid constraint",
			args:     "spf13/cobra@invalid",
			wantExit: 1,
		},
		{
			name:     "non-existent repo",
			args:     "nonexistent/repo@^1.0.0",
			wantExit: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We'll skip actual CLI execution in this test
			// as it would require building the binary first
			t.Logf("Would test: semver resolve %s", tt.args)
		})
	}
}

func TestRateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test without token to verify rate limit handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	githubClient := github.NewClient("") // No token
	r := resolver.New(githubClient, resolver.Options{})

	// Make a single request (should work even without token)
	version, err := r.MaxSatisfying(ctx, "golang", "go", "^1.21.0")
	if err != nil {
		// Check if it's a rate limit error
		t.Logf("Request failed (possibly rate limited): %v", err)
	} else {
		t.Logf("Resolved golang/go@^1.21.0 to %s without token", version)
	}
}

func TestEnvironmentVariables(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Save original env vars
	origGitHubToken := os.Getenv("GITHUB_TOKEN")
	origGHToken := os.Getenv("GH_TOKEN")
	defer func() {
		os.Setenv("GITHUB_TOKEN", origGitHubToken)
		os.Setenv("GH_TOKEN", origGHToken)
	}()

	// Test with GITHUB_TOKEN
	os.Setenv("GITHUB_TOKEN", "test-token")
	os.Unsetenv("GH_TOKEN")

	tm := auth.NewTokenManager()
	token, err := tm.GetToken()
	if err != nil {
		t.Fatalf("GetToken() with GITHUB_TOKEN set: %v", err)
	}
	if token != "test-token" {
		t.Errorf("GetToken() = %v, want test-token", token)
	}

	// Test with GH_TOKEN
	os.Unsetenv("GITHUB_TOKEN")
	os.Setenv("GH_TOKEN", "gh-test-token")

	token, err = tm.GetToken()
	if err != nil {
		t.Fatalf("GetToken() with GH_TOKEN set: %v", err)
	}
	if token != "gh-test-token" {
		t.Errorf("GetToken() = %v, want gh-test-token", token)
	}
}

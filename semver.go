// Package semver provides a library for resolving semantic version constraints
// against GitHub repository tags and releases.
package semver

import (
	"context"

	"github.com/binary-install/semver/pkg/auth"
	"github.com/binary-install/semver/pkg/github"
	"github.com/binary-install/semver/pkg/resolver"
)

// Options configures the version resolution behavior.
type Options struct {
	// Token is the GitHub personal access token.
	// If empty, it will try to get token from environment variables or keyring.
	Token string

	// IncludePrerelease includes prerelease versions in the resolution.
	IncludePrerelease bool

	// IncludeDraft includes draft releases in the resolution.
	IncludeDraft bool

	// IncludeTags includes git tags in addition to releases.
	// By default, only GitHub releases are considered.
	IncludeTags bool

	// GitHubClient allows using a custom GitHub client.
	// If nil, a default client will be created.
	GitHubClient github.Client
}

// MaxSatisfying returns the highest version that satisfies the given constraint
// for the specified GitHub repository.
//
// The constraint should be a valid semantic version constraint as defined by
// https://github.com/Masterminds/semver.
//
// Example:
//
//	version, err := semver.MaxSatisfying(ctx, "golang", "go", "~1.21.0", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(version) // e.g., "go1.21.5"
func MaxSatisfying(ctx context.Context, owner, repo, constraint string, opts *Options) (string, error) {
	if opts == nil {
		opts = &Options{}
	}

	// Get GitHub client
	var githubClient github.Client
	if opts.GitHubClient != nil {
		githubClient = opts.GitHubClient
	} else {
		token := opts.Token
		if token == "" {
			// Try to get token from environment or keyring
			tm := auth.NewTokenManager()
			if t, err := tm.GetToken(); err == nil {
				token = t
			}
			// If no token found, continue with anonymous client
		}
		githubClient = github.NewClient(token)
	}

	// Create resolver
	resolverOpts := resolver.Options{
		IncludePrerelease: opts.IncludePrerelease,
		IncludeDraft:      opts.IncludeDraft,
		IncludeTags:       opts.IncludeTags,
	}
	r := resolver.New(githubClient, resolverOpts)

	// Resolve version
	return r.MaxSatisfying(ctx, owner, repo, constraint)
}

// MaxSatisfyingWithToken is a convenience function that creates a GitHub client
// with the provided token and resolves the version constraint.
//
// This is equivalent to calling MaxSatisfying with Options{Token: token}.
func MaxSatisfyingWithToken(ctx context.Context, owner, repo, constraint, token string) (string, error) {
	return MaxSatisfying(ctx, owner, repo, constraint, &Options{Token: token})
}

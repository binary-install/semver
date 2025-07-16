package github

import (
	"context"
)

// Client provides access to GitHub repository information.
type Client interface {
	// ListTags returns all tags for the specified repository.
	ListTags(ctx context.Context, owner, repo string) ([]string, error)
	// ListReleases returns all releases for the specified repository.
	ListReleases(ctx context.Context, owner, repo string) ([]Release, error)
}

// Release represents a GitHub release.
type Release struct {
	Tag        string
	Prerelease bool
	Draft      bool
}

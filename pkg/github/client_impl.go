package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
	semvererrors "github.com/binary-install/semver/internal/errors"
)

const (
	// Default per page for pagination
	defaultPerPage = 100
)

// clientImpl implements the Client interface.
type clientImpl struct {
	client *github.Client
}

// NewClient creates a new GitHub client.
func NewClient(token string) Client {
	var httpClient *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient = oauth2.NewClient(context.Background(), ts)
	}

	return &clientImpl{
		client: github.NewClient(httpClient),
	}
}

// ListTags returns all tags for the specified repository.
func (c *clientImpl) ListTags(ctx context.Context, owner, repo string) ([]string, error) {
	if owner == "" || repo == "" {
		return nil, semvererrors.NewValidationError("owner and repo cannot be empty")
	}

	var allTags []string
	opts := &github.ListOptions{
		PerPage: defaultPerPage,
	}

	for {
		tags, resp, err := c.client.Repositories.ListTags(ctx, owner, repo, opts)
		if err != nil {
			return nil, c.wrapGitHubError(err, owner, repo)
		}

		for _, tag := range tags {
			if tag.Name != nil {
				allTags = append(allTags, *tag.Name)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

		// Add a small delay to be respectful of rate limits
		select {
		case <-ctx.Done():
			return nil, semvererrors.WrapError(ctx.Err(), semvererrors.ErrorTypeNetwork, "context cancelled")
		case <-time.After(100 * time.Millisecond):
		}
	}

	return allTags, nil
}

// ListReleases returns all releases for the specified repository.
func (c *clientImpl) ListReleases(ctx context.Context, owner, repo string) ([]Release, error) {
	if owner == "" || repo == "" {
		return nil, semvererrors.NewValidationError("owner and repo cannot be empty")
	}

	var allReleases []Release
	opts := &github.ListOptions{
		PerPage: defaultPerPage,
	}

	for {
		releases, resp, err := c.client.Repositories.ListReleases(ctx, owner, repo, opts)
		if err != nil {
			return nil, c.wrapGitHubError(err, owner, repo)
		}

		for _, release := range releases {
			if release.TagName != nil {
				allReleases = append(allReleases, Release{
					Tag:        *release.TagName,
					Prerelease: release.Prerelease != nil && *release.Prerelease,
					Draft:      release.Draft != nil && *release.Draft,
				})
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

		// Add a small delay to be respectful of rate limits
		select {
		case <-ctx.Done():
			return nil, semvererrors.WrapError(ctx.Err(), semvererrors.ErrorTypeNetwork, "context cancelled")
		case <-time.After(100 * time.Millisecond):
		}
	}

	return allReleases, nil
}

// wrapGitHubError wraps GitHub API errors with appropriate error types.
func (c *clientImpl) wrapGitHubError(err error, owner, repo string) error {
	if err == nil {
		return nil
	}

	if ghErr, ok := err.(*github.ErrorResponse); ok {
		switch ghErr.Response.StatusCode {
		case http.StatusNotFound:
			return semvererrors.NewNotFoundError(fmt.Sprintf("repository %s/%s not found", owner, repo))
		case http.StatusUnauthorized:
			return semvererrors.NewAuthenticationError("invalid or expired GitHub token")
		case http.StatusForbidden:
			if ghErr.Message == "API rate limit exceeded" {
				return semvererrors.NewRateLimitError("GitHub API rate limit exceeded")
			}
			return semvererrors.NewAuthenticationError("insufficient permissions to access repository")
		case http.StatusTooManyRequests:
			return semvererrors.NewRateLimitError("GitHub API rate limit exceeded")
		default:
			return semvererrors.WrapError(err, semvererrors.ErrorTypeNetwork, fmt.Sprintf("GitHub API error: %s", ghErr.Message))
		}
	}

	return semvererrors.WrapError(err, semvererrors.ErrorTypeNetwork, "failed to communicate with GitHub API")
}

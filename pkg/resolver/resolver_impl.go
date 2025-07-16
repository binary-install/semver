package resolver

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/binary-install/semver/pkg/github"
	semvererrors "github.com/binary-install/semver/internal/errors"
)

// resolverImpl implements the Resolver interface.
type resolverImpl struct {
	githubClient github.Client
	options      Options
}

// New creates a new resolver instance.
func New(githubClient github.Client, options Options) Resolver {
	return &resolverImpl{
		githubClient: githubClient,
		options:      options,
	}
}

// MaxSatisfying returns the highest version that satisfies the given constraint.
func (r *resolverImpl) MaxSatisfying(ctx context.Context, owner, repo, constraint string) (string, error) {
	if owner == "" || repo == "" || constraint == "" {
		return "", semvererrors.NewValidationError("owner, repo, and constraint cannot be empty")
	}

	// Parse the constraint
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return "", semvererrors.WrapError(err, semvererrors.ErrorTypeParse, fmt.Sprintf("invalid constraint: %s", constraint))
	}

	// Get all versions from tags and releases
	versions, err := r.getAllVersions(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", semvererrors.NewNotFoundError(fmt.Sprintf("no versions found for %s/%s", owner, repo))
	}

	// Find the maximum satisfying version
	var maxVersion *semver.Version
	var maxVersionString string

	for versionString, version := range versions {
		if c.Check(version) {
			if maxVersion == nil || version.GreaterThan(maxVersion) {
				maxVersion = version
				maxVersionString = versionString
			}
		}
	}

	if maxVersion == nil {
		return "", semvererrors.NewNotFoundError(fmt.Sprintf("no version satisfies constraint %s", constraint))
	}

	return maxVersionString, nil
}

// getAllVersions fetches all versions from both tags and releases.
func (r *resolverImpl) getAllVersions(ctx context.Context, owner, repo string) (map[string]*semver.Version, error) {
	versions := make(map[string]*semver.Version)

	// Fetch releases first (they have more metadata)
	releases, err := r.githubClient.ListReleases(ctx, owner, repo)
	if err != nil {
		// If releases fail, we'll try tags
		releases = []github.Release{}
	}

	// Add release versions
	for _, release := range releases {
		// Skip drafts unless explicitly included
		if release.Draft && !r.options.IncludeDraft {
			continue
		}
		// Skip prereleases unless explicitly included
		if release.Prerelease && !r.options.IncludePrerelease {
			continue
		}

		version, err := parseVersion(release.Tag)
		if err != nil {
			continue // Skip invalid versions
		}

		versions[release.Tag] = version
	}

	// Fetch tags
	tags, err := r.githubClient.ListTags(ctx, owner, repo)
	if err != nil {
		// If both releases and tags fail, return error
		if len(versions) == 0 {
			return nil, err
		}
		// Otherwise, continue with what we have from releases
		return versions, nil
	}

	// Add tag versions (only if not already present from releases)
	for _, tag := range tags {
		if _, exists := versions[tag]; exists {
			continue
		}

		version, err := parseVersion(tag)
		if err != nil {
			continue // Skip invalid versions
		}

		// For tags, we can't determine if it's a prerelease from metadata,
		// so we check the version itself
		if version.Prerelease() != "" && !r.options.IncludePrerelease {
			continue
		}

		versions[tag] = version
	}

	return versions, nil
}

// parseVersion parses a version string, handling common prefixes.
func parseVersion(versionStr string) (*semver.Version, error) {
	// Remove common prefixes
	cleaned := versionStr
	cleaned = strings.TrimPrefix(cleaned, "v")
	cleaned = strings.TrimPrefix(cleaned, "V")
	cleaned = strings.TrimPrefix(cleaned, "go") // For Go versions like "go1.21.0"
	cleaned = strings.TrimPrefix(cleaned, "release-")
	cleaned = strings.TrimPrefix(cleaned, "release/")

	// Try to parse the cleaned version
	version, err := semver.NewVersion(cleaned)
	if err != nil {
		// Try the original string as a last resort
		version, err = semver.NewVersion(versionStr)
		if err != nil {
			return nil, err
		}
	}

	return version, nil
}

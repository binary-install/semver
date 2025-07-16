package resolver

import (
	"context"
)

// Resolver resolves semantic version constraints against repository versions.
type Resolver interface {
	// MaxSatisfying returns the highest version that satisfies the given constraint.
	// Returns an empty string if no version satisfies the constraint.
	MaxSatisfying(ctx context.Context, owner, repo, constraint string) (string, error)
}

// Options configures the resolver behavior.
type Options struct {
	// IncludePrerelease includes prerelease versions in the resolution.
	IncludePrerelease bool
	// IncludeDraft includes draft releases in the resolution.
	IncludeDraft bool
	// IncludeTags includes git tags in addition to releases.
	// By default, only GitHub releases are considered.
	IncludeTags bool
}

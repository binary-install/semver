package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/binary-install/semver/pkg/auth"
	"github.com/binary-install/semver/pkg/github"
	"github.com/binary-install/semver/pkg/resolver"
	semvererrors "github.com/binary-install/semver/internal/errors"
)

var (
	token             string
	includePrerelease bool
	includeDraft      bool
	includeTags       bool
	timeout           time.Duration
)

// resolveCmd represents the resolve command
var resolveCmd = &cobra.Command{
	Use:   "resolve <owner>/<repo>@<constraint>",
	Short: "Resolve semantic version constraint",
	Long: `Resolve finds the highest version that satisfies the given constraint.

Examples:
  # Find exact version
  semver resolve golang/go@1.21.0

  # Find latest patch version
  semver resolve golang/go@~1.21.0

  # Find latest minor version
  semver resolve golang/go@^1.21.0

  # Find version within range
  semver resolve golang/go@">=1.20.0 <1.22.0"

  # Include prereleases
  semver resolve --prerelease golang/go@^1.21.0`,
	Args: cobra.ExactArgs(1),
	RunE: runResolve,
}

func init() {
	rootCmd.AddCommand(resolveCmd)

	resolveCmd.Flags().StringVar(&token, "token", "", "GitHub personal access token (defaults to GITHUB_TOKEN or GH_TOKEN env var)")
	resolveCmd.Flags().BoolVar(&includePrerelease, "prerelease", false, "Include prerelease versions")
	resolveCmd.Flags().BoolVar(&includeDraft, "draft", false, "Include draft releases")
	resolveCmd.Flags().BoolVar(&includeTags, "tags", false, "Include git tags in addition to releases")
	resolveCmd.Flags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")
}

func runResolve(cmd *cobra.Command, args []string) error {
	// Parse the input format: owner/repo@constraint
	input := args[0]
	parts := strings.SplitN(input, "@", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format: expected <owner>/<repo>@<constraint>, got %s", input)
	}

	repoParts := strings.SplitN(parts[0], "/", 2)
	if len(repoParts) != 2 {
		return fmt.Errorf("invalid repository format: expected <owner>/<repo>, got %s", parts[0])
	}

	owner := repoParts[0]
	repo := repoParts[1]
	constraint := parts[1]

	// Get token
	if token == "" {
		tm := auth.NewTokenManager()
		t, err := tm.GetToken()
		if err != nil {
			// Continue without token (unauthenticated requests have lower rate limits)
			cmd.PrintErrln("Warning: No GitHub token found. API rate limits will be restrictive (60 requests/hour).")
			cmd.PrintErrln("To increase rate limits to 5,000 requests/hour, you can either:")
			cmd.PrintErrln("  1. Set GITHUB_TOKEN environment variable: export GITHUB_TOKEN=ghp_...")
			cmd.PrintErrln("  2. Use semver token management: semver token set")
		} else {
			token = t
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create GitHub client
	githubClient := github.NewClient(token)

	// Create resolver
	opts := resolver.Options{
		IncludePrerelease: includePrerelease,
		IncludeDraft:      includeDraft,
		IncludeTags:       includeTags,
	}
	r := resolver.New(githubClient, opts)

	// Resolve version
	version, err := r.MaxSatisfying(ctx, owner, repo, constraint)
	if err != nil {
		return handleError(cmd, err)
	}

	// Output the result
	fmt.Println(version)
	return nil
}

func handleError(cmd *cobra.Command, err error) error {
	var typedErr semvererrors.TypedError
	if errors, ok := err.(semvererrors.TypedError); ok {
		typedErr = errors
	}

	if typedErr != nil {
		switch typedErr.Type() {
		case semvererrors.ErrorTypeAuthentication:
			cmd.PrintErrln("Error: Authentication failed. Please check your GitHub token.")
			return err
		case semvererrors.ErrorTypeRateLimit:
			cmd.PrintErrln("Error: GitHub API rate limit exceeded. Please try again later or authenticate with a token.")
			return err
		case semvererrors.ErrorTypeNotFound:
			cmd.PrintErrln(fmt.Sprintf("Error: %v", err))
			return err
		case semvererrors.ErrorTypeParse:
			cmd.PrintErrln(fmt.Sprintf("Error: Invalid version constraint: %v", err))
			return err
		case semvererrors.ErrorTypeNetwork:
			cmd.PrintErrln(fmt.Sprintf("Error: Network error: %v", err))
			return err
		case semvererrors.ErrorTypeValidation:
			cmd.PrintErrln(fmt.Sprintf("Error: Validation error: %v", err))
			return err
		}
	}

	// Generic error
	cmd.PrintErrln(fmt.Sprintf("Error: %v", err))
	return err
}

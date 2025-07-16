package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Version is set via build flags
	Version = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "semver",
	Short: "Resolve semantic version constraints against GitHub repository tags",
	Long: `semver is a tool for resolving semantic version constraints against
GitHub repository tags and releases. It finds the highest version that
matches a given semantic version constraint.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

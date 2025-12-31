package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/cmd"
	"github.com/velocitykode/velocity-cli/internal/delegator"
	"github.com/velocitykode/velocity-cli/internal/version"
)

func main() {
	// Check Go version immediately - Velocity requires Go 1.25+
	if err := version.CheckGoVersion(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Check if we should delegate to project CLI
	// This happens when:
	// 1. We're in a Velocity project (has cmd/velocity/main.go)
	// 2. The command is not a global-only command (new, init, help, etc.)
	if delegator.ShouldDelegate(os.Args[1:]) {
		// Delegate to project's CLI
		if err := delegator.Delegate(os.Args[1:]); err != nil {
			os.Exit(1)
		}
		return
	}

	// Run global CLI
	rootCmd := &cobra.Command{
		Use:     "velocity",
		Short:   "CLI for the Velocity Go web framework",
		Version: cmd.Version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Global commands (always available)
	rootCmd.AddCommand(cmd.NewCmd)
	rootCmd.AddCommand(cmd.InitCmd)
	rootCmd.AddCommand(cmd.ConfigCmd)
	rootCmd.AddCommand(cmd.VersionCmd)

	// Initialize help after adding all commands
	cmd.InitHelp(rootCmd)

	// Execute with proper exit code
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

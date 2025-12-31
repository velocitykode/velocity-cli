// Package cli provides an importable CLI for Velocity projects.
// User projects import this package in cmd/velocity/main.go to get
// access to commands like serve, migrate, make:controller, etc.
package cli

import (
	"github.com/spf13/cobra"
)

// Version is the CLI version
var Version = "0.5.0"

var rootCmd *cobra.Command

func init() {
	initRootCmd()
}

func initRootCmd() {
	rootCmd = &cobra.Command{
		Use:     "velocity",
		Short:   "Velocity CLI - Development tools for Velocity projects",
		Version: Version,
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Register all commands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateFreshCmd)
	rootCmd.AddCommand(makeControllerCmd)
	rootCmd.AddCommand(keyGenerateCmd)
}

// Execute runs the CLI
func Execute() error {
	if rootCmd == nil {
		initRootCmd()
	}
	return rootCmd.Execute()
}

// RootCmd returns the root command for testing
func RootCmd() *cobra.Command {
	return rootCmd
}

// trigger

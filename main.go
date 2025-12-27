package main

import (
	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/cmd"
	"github.com/velocitykode/velocity-cli/framework"
	"os"
)

func main() {
	// Create unified root command
	rootCmd := &cobra.Command{
		Use:     "velocity",
		Short:   "CLI for the Velocity Go web framework",
		Version: "0.1.0",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Project commands
	rootCmd.AddCommand(cmd.NewCmd)
	rootCmd.AddCommand(cmd.InitCmd)
	rootCmd.AddCommand(cmd.ConfigCmd)

	// Development commands
	rootCmd.AddCommand(framework.ServeCmd)
	rootCmd.AddCommand(framework.BuildCmd)

	// Database commands
	rootCmd.AddCommand(framework.MigrateCmd)
	rootCmd.AddCommand(framework.MigrateFreshCmd)

	// Code generation
	makeCmd := &cobra.Command{
		Use:     "make",
		Aliases: []string{"generate", "g"},
		Short:   "Generate code (controllers, models, etc.)",
	}
	makeCmd.AddCommand(framework.MakeControllerCmd)
	rootCmd.AddCommand(makeCmd)

	// Utility commands
	rootCmd.AddCommand(cmd.KeyCmd)

	// Info commands
	rootCmd.AddCommand(cmd.VersionCmd)

	// Initialize help after adding all commands
	cmd.InitHelp(rootCmd)

	// Execute with proper exit code
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

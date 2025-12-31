package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/generator"
	"github.com/velocitykode/velocity-cli/internal/ui"
)

var (
	database string
	cache    string
	auth     bool
	api      bool
)

var NewCmd = &cobra.Command{
	Use:           "new [project-name]",
	Short:         "Create a new Velocity project",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			ui.Error("Project name is required")
			ui.Newline()
			ui.Muted("Usage: velocity new [project-name]")
			ui.Newline()
			ui.Muted("Flags:")
			ui.Muted("  --database    Database driver (postgres, mysql, sqlite)")
			ui.Muted("  --cache       Cache driver (redis, memory)")
			ui.Muted("  --auth        Include authentication scaffolding")
			ui.Muted("  --api         API-only structure (no views)")
			return fmt.Errorf("")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		ui.Header("velocity new")

		// Create project with flags (defaults to sqlite if not specified)
		config := generator.ProjectConfig{
			Name:     projectName,
			Module:   projectName,
			Database: database,
			Cache:    cache,
			Auth:     auth,
			API:      api,
		}

		if err := generator.CreateProject(config); err != nil {
			ui.Newline()
			ui.Error(err.Error())
			return
		}

		ui.Newline()
		ui.Info("Starting development servers")

		generator.StartDevServers(projectName)
	},
}

func init() {
	NewCmd.Flags().StringVar(&database, "database", "sqlite", "Database driver (postgres, mysql, sqlite)")
	NewCmd.Flags().StringVar(&cache, "cache", "memory", "Cache driver (redis, memory)")
	NewCmd.Flags().BoolVar(&auth, "auth", false, "Include authentication scaffolding")
	NewCmd.Flags().BoolVar(&api, "api", false, "API-only structure (no views)")
}

package cmd

import (
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
	Use:   "new [project-name]",
	Short: "Create a new Velocity project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

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
			ui.Error("Error creating project: " + err.Error())
			return
		}

		ui.Newline()
		ui.Success("Project created successfully!")
		ui.Newline()
		ui.Info("Starting development servers")
		ui.Newline()

		generator.StartDevServers(projectName)
	},
}

func init() {
	NewCmd.Flags().StringVar(&database, "database", "sqlite", "Database driver (postgres, mysql, sqlite)")
	NewCmd.Flags().StringVar(&cache, "cache", "memory", "Cache driver (redis, memory)")
	NewCmd.Flags().BoolVar(&auth, "auth", false, "Include authentication scaffolding")
	NewCmd.Flags().BoolVar(&api, "api", false, "API-only structure (no views)")
}

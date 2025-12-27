package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/generator"
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
			fmt.Printf("\nError creating project: %v\n", err)
			return
		}

		// Show success message
		const (
			colorReset = "\033[0m"
			colorGreen = "\033[32m"
			colorCyan  = "\033[36m"
			colorBlue  = "\033[34m"
			bold       = "\033[1m"
			underline  = "\033[4m"
		)

		fmt.Println()
		fmt.Printf("\n%s  %s\n", bold+underline+colorGreen+"SUCCESS"+colorReset, "Project created successfully!")
		fmt.Println()

		// Start npm dev server in background
		fmt.Printf("%s  %s\n\n", bold+underline+colorBlue+"INFO"+colorReset, "Starting development servers")
		fmt.Println()

		generator.StartDevServers(projectName)
	},
}

func init() {
	NewCmd.Flags().StringVar(&database, "database", "sqlite", "Database driver (postgres, mysql, sqlite)")
	NewCmd.Flags().StringVar(&cache, "cache", "memory", "Cache driver (redis, memory)")
	NewCmd.Flags().BoolVar(&auth, "auth", false, "Include authentication scaffolding")
	NewCmd.Flags().BoolVar(&api, "api", false, "API-only structure (no views)")
}

package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/delegator"
	"github.com/velocitykode/velocity-cli/internal/detector"
	"github.com/velocitykode/velocity-cli/internal/ui"
)

// UpgradeCmd updates the project's velocity-cli dependency to match the global CLI version
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Update project's velocity-cli to match global CLI version",
	Long:  `Updates the velocity-cli dependency in your project's go.mod to match the globally installed CLI version.`,
	RunE:  runUpgrade,
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	// Check if we're in a Velocity project
	if !detector.IsVelocityProject() {
		return fmt.Errorf("not in a Velocity project (no cmd/velocity/main.go found)")
	}

	ui.Header("upgrade")

	// Get versions
	globalVersion := Version
	projectVersion := delegator.GetProjectCLIVersion()

	if projectVersion == "" {
		return fmt.Errorf("could not determine project's velocity-cli version from go.mod")
	}

	// Normalize versions
	global := strings.TrimPrefix(globalVersion, "v")
	project := strings.TrimPrefix(projectVersion, "v")

	if global == project {
		ui.Success(fmt.Sprintf("Already up to date (v%s)", global))
		return nil
	}

	ui.Info(fmt.Sprintf("Upgrading velocity-cli: %s -> v%s", projectVersion, global))

	// Run go get to update
	goGet := exec.Command("go", "get", fmt.Sprintf("github.com/velocitykode/velocity-cli@v%s", global))
	goGet.Stdout = nil
	goGet.Stderr = nil

	if err := goGet.Run(); err != nil {
		ui.Error("Failed to update velocity-cli")
		return err
	}

	// Run go mod tidy
	goTidy := exec.Command("go", "mod", "tidy")
	goTidy.Stdout = nil
	goTidy.Stderr = nil
	goTidy.Run() // Ignore error, tidy is optional

	// Clear cached binary so it rebuilds on next command
	exec.Command("rm", "-f", ".velocity/bin/cli").Run()

	ui.Success(fmt.Sprintf("Updated to v%s", global))
	ui.Muted("Project CLI will rebuild on next command")

	return nil
}

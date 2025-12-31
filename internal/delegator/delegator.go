// Package delegator handles delegation from global CLI to project CLI.
// When running commands inside a Velocity project, the global CLI
// delegates to `go run ./cmd/velocity` (or cached binary) so the
// project's CLI has access to migrations and other project-specific code.
package delegator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/velocitykode/velocity-cli/internal/detector"
	"github.com/velocitykode/velocity-cli/internal/ui"
)

// GlobalCommands are commands that should NOT be delegated.
// These run in the global CLI context.
var GlobalCommands = map[string]bool{
	"new":         true,
	"init":        true,
	"upgrade":     true,
	"self-update": true,
	"help":        true,
	"--help":      true,
	"-h":          true,
	"version":     true,
	"--version":   true,
	"-v":          true,
	"config":      true,
}

// ShouldDelegate returns true if the command should be delegated
// to the project's CLI (cmd/velocity/main.go).
func ShouldDelegate(args []string) bool {
	if len(args) == 0 {
		return false
	}

	// Check if it's a global command
	if GlobalCommands[args[0]] {
		return false
	}

	// Check if we're in a Velocity project
	return detector.IsVelocityProject()
}

// Delegate runs the command via the project's CLI.
// It first checks for a cached binary, rebuilding if necessary.
func Delegate(args []string) error {
	// Check for cached binary
	cachedBin := ".velocity/bin/cli"

	if needsRebuild(cachedBin) {
		ui.Step("Building project CLI...")

		// Ensure directory exists
		os.MkdirAll(filepath.Dir(cachedBin), 0755)

		// Build the project CLI
		buildCmd := exec.Command("go", "build", "-o", cachedBin, "./cmd/velocity")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		if err := buildCmd.Run(); err != nil {
			// Fall back to go run if build fails
			return runWithGoRun(args)
		}
	}

	// Run cached binary
	cmd := exec.Command(cachedBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// runWithGoRun falls back to using go run if binary caching fails.
func runWithGoRun(args []string) error {
	cmdArgs := append([]string{"run", "./cmd/velocity"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// needsRebuild checks if the cached binary needs to be rebuilt.
// Returns true if:
// - Binary doesn't exist
// - go.mod or go.sum is newer than binary
// - Any .go file in cmd/velocity is newer than binary
func needsRebuild(binPath string) bool {
	binInfo, err := os.Stat(binPath)
	if err != nil {
		return true // Binary doesn't exist
	}

	binTime := binInfo.ModTime()

	// Check go.mod
	if isNewer("go.mod", binTime) {
		return true
	}

	// Check go.sum
	if isNewer("go.sum", binTime) {
		return true
	}

	// Check cmd/velocity directory
	cmdDir := "cmd/velocity"
	if _, err := os.Stat(cmdDir); err != nil {
		return true // cmd/velocity doesn't exist
	}

	// Walk cmd/velocity and check if any .go file is newer
	err = filepath.Walk(cmdDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if info.ModTime().After(binTime) {
				return fmt.Errorf("rebuild needed")
			}
		}
		return nil
	})

	return err != nil
}

// isNewer checks if a file is newer than the given time.
func isNewer(path string, t time.Time) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.ModTime().After(t)
}

// CheckVersionMismatch warns if global CLI version differs from project's cli package version.
func CheckVersionMismatch(globalVersion string) {
	projectVersion := getProjectCLIVersion()
	if projectVersion != "" && projectVersion != globalVersion {
		ui.Warning(fmt.Sprintf("CLI version mismatch: global=%s, project=%s", globalVersion, projectVersion))
	}
}

// getProjectCLIVersion extracts the velocity-cli version from go.mod.
func getProjectCLIVersion() string {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "github.com/velocitykode/velocity-cli") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}

	return ""
}

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

// Optional directories that affect the CLI binary when changed.
// These are imported by cmd/velocity/main.go.
var cliOptionalDirs = []string{
	"database/migrations",
	"bootstrap",
}

// needsRebuild checks if the cached binary needs to be rebuilt.
// Returns true if:
// - Binary doesn't exist
// - go.mod or go.sum is newer than binary
// - Any .go file in CLI source directories is newer than binary
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

	// cmd/velocity is required - if missing, rebuild needed
	if _, err := os.Stat("cmd/velocity"); err != nil {
		return true
	}

	// Check cmd/velocity (required)
	if hasNewerGoFiles("cmd/velocity", binTime) {
		return true
	}

	// Check optional directories that CLI imports
	for _, dir := range cliOptionalDirs {
		if hasNewerGoFiles(dir, binTime) {
			return true
		}
	}

	return false
}

// hasNewerGoFiles checks if any .go file in dir is newer than the given time.
// Returns true if newer files found OR if there's an error reading the directory
// (fail-safe: rebuild if we can't verify).
func hasNewerGoFiles(dir string, t time.Time) bool {
	if _, err := os.Stat(dir); err != nil {
		return false // Directory doesn't exist
	}

	var found bool
	var walkErr bool
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			walkErr = true
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if info.ModTime().After(t) {
				found = true
				return filepath.SkipAll
			}
		}
		return nil
	})
	return found || walkErr
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
	if projectVersion == "" {
		return
	}

	// Normalize: strip "v" prefix for comparison
	global := strings.TrimPrefix(globalVersion, "v")
	project := strings.TrimPrefix(projectVersion, "v")

	if global != project {
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

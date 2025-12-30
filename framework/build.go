package framework

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/ui"
)

var (
	output   string
	targetOS string
	arch     string
	optimize bool
	version  string
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the application for production",
	Long: `Build compiles your Velocity application for production deployment.
	
It supports cross-compilation, optimization, and versioning.`,
	Run: func(cmd *cobra.Command, args []string) {
		buildApp()
	},
}

func init() {

	BuildCmd.Flags().StringVarP(&output, "output", "o", "./dist/app", "Output path for the binary")
	BuildCmd.Flags().StringVar(&targetOS, "os", runtime.GOOS, "Target operating system")
	BuildCmd.Flags().StringVar(&arch, "arch", runtime.GOARCH, "Target architecture")
	BuildCmd.Flags().BoolVar(&optimize, "optimize", true, "Enable optimizations")
	BuildCmd.Flags().StringVar(&version, "version", "", "Version to embed in binary")
}

func buildApp() {
	ui.Header("build")
	ui.Step("Building Velocity application...")

	// Prepare build arguments
	args := []string{"build"}

	// Output file
	args = append(args, "-o", output)

	// Optimizations
	if optimize {
		args = append(args, "-ldflags", "-s -w")
		ui.Step("Building with optimizations...")
	}

	// Version info
	if version != "" {
		ldflags := fmt.Sprintf("-X main.Version=%s", version)
		if optimize {
			ldflags = "-s -w " + ldflags
		}
		args = append(args, "-ldflags", ldflags)
		ui.KeyValue("Version", version)
	}

	// Add the main package
	args = append(args, ".")

	// Set environment for cross-compilation
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", targetOS),
		fmt.Sprintf("GOARCH=%s", arch),
		"CGO_ENABLED=0", // Disable CGO for better portability
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	ui.KeyValue("Target", fmt.Sprintf("%s/%s", targetOS, arch))

	// Run the build
	if err := cmd.Run(); err != nil {
		ui.Error(fmt.Sprintf("Build failed: %v", err))
		os.Exit(1)
	}

	// Get file info
	info, err := os.Stat(output)
	if err == nil {
		size := info.Size()
		sizeStr := formatBytes(size)
		ui.Success(fmt.Sprintf("Built successfully: %s (%s)", output, sizeStr))

		// Make executable on Unix systems
		if targetOS != "windows" {
			os.Chmod(output, 0755)
		}

		// Show next steps
		steps := []string{
			fmt.Sprintf("Deploy %s to your server", output),
			"Set environment variables",
			"Run the binary",
		}
		if targetOS == runtime.GOOS && arch == runtime.GOARCH {
			steps = append(steps, fmt.Sprintf("Test locally: %s", output))
		}
		ui.NextSteps(steps)
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

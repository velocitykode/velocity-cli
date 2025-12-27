package framework

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
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
	fmt.Println("🔨 Building Velocity application...")

	// Prepare build arguments
	args := []string{"build"}

	// Output file
	args = append(args, "-o", output)

	// Optimizations
	if optimize {
		args = append(args, "-ldflags", "-s -w")
		fmt.Println("📦 Building with optimizations...")
	}

	// Version info
	if version != "" {
		ldflags := fmt.Sprintf("-X main.Version=%s", version)
		if optimize {
			ldflags = "-s -w " + ldflags
		}
		args = append(args, "-ldflags", ldflags)
		fmt.Printf("🏷️  Version: %s\n", version)
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

	fmt.Printf("Target: %s/%s\n", targetOS, arch)

	// Run the build
	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}

	// Get file info
	info, err := os.Stat(output)
	if err == nil {
		size := info.Size()
		sizeStr := formatBytes(size)
		fmt.Printf("Built successfully: %s (%s)\n", output, sizeStr)

		// Make executable on Unix systems
		if targetOS != "windows" {
			os.Chmod(output, 0755)
		}

		// Show next steps
		fmt.Println("\n📋 Next steps:")
		fmt.Printf("   1. Deploy %s to your server\n", output)
		fmt.Println("   2. Set environment variables")
		fmt.Println("   3. Run the binary")

		if targetOS == runtime.GOOS && arch == runtime.GOARCH {
			fmt.Printf("\n   Test locally: %s\n", output)
		}
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

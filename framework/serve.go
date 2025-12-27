package framework

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/colors"
)

var (
	port      string
	env       string
	watch     bool
	buildTags string
)

// ServeCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the development server",
	Long: `Start the Velocity development server with optional hot reload.
	
The server will automatically reload when Go files change if --watch is enabled.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		fmt.Println(colors.BrandStyle.Bold(true).Render("VELOCITY CLI SERVER"))
		if watch {
			fmt.Printf("Starting on port %s with hot reload...\n", colors.SuccessStyle.Render(port))
		} else {
			fmt.Printf("Starting on port %s...\n", colors.SuccessStyle.Render(port))
		}
		fmt.Println()

		if watch {
			runWithWatcher()
		} else {
			runServer()
		}
	},
}

func init() {

	ServeCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server on")
	ServeCmd.Flags().StringVarP(&env, "env", "e", "development", "Environment to run in")
	ServeCmd.Flags().BoolVarP(&watch, "watch", "w", true, "Enable hot reload")
	ServeCmd.Flags().StringVar(&buildTags, "tags", "", "Build tags to pass to go build")
}

func runServer() {
	// Set environment variables
	os.Setenv("APP_ENV", env)
	os.Setenv("APP_PORT", port)

	// Build and run the application
	buildCmd := exec.Command("go", "build", "-o", ".velocity/tmp/server", ".")
	if buildTags != "" {
		buildCmd.Args = append(buildCmd.Args[:2], "-tags", buildTags)
		buildCmd.Args = append(buildCmd.Args, buildCmd.Args[2:]...)
	}

	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	// Run the server
	serverCmd := exec.Command(".velocity/tmp/server")
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr
	serverCmd.Env = os.Environ()

	if err := serverCmd.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func runWithWatcher() {
	// Create temp directory for builds
	os.MkdirAll(".velocity/tmp", 0755)

	// Channel to trigger rebuilds
	rebuild := make(chan bool, 1)

	// Start file watcher
	go watchFiles(rebuild)

	// Server management
	var serverCmd *exec.Cmd
	var mu sync.Mutex

	// Function to start/restart server
	startServer := func() {
		mu.Lock()
		defer mu.Unlock()

		// Kill existing server if running
		if serverCmd != nil && serverCmd.Process != nil {
			fmt.Println("Stopping server...")
			serverCmd.Process.Kill()
			serverCmd.Wait()
		}

		// Build
		fmt.Println("Building...")
		buildCmd := exec.Command("go", "build", "-o", ".velocity/tmp/server", ".")
		if buildTags != "" {
			buildCmd.Args = append(buildCmd.Args[:2], "-tags", buildTags)
			buildCmd.Args = append(buildCmd.Args, buildCmd.Args[2:]...)
		}

		if output, err := buildCmd.CombinedOutput(); err != nil {
			fmt.Printf("%s\n%s\n", colors.ErrorStyle.Render("Build failed:"), output)
			return
		}

		// Start new server
		fmt.Printf("%s on port %s...\n", colors.SuccessStyle.Render("Starting server"), port)
		serverCmd = exec.Command(".velocity/tmp/server")
		serverCmd.Stdout = os.Stdout
		serverCmd.Stderr = os.Stderr
		serverCmd.Env = append(os.Environ(),
			fmt.Sprintf("APP_ENV=%s", env),
			fmt.Sprintf("APP_PORT=%s", port),
		)

		if err := serverCmd.Start(); err != nil {
			fmt.Printf("❌ Failed to start server: %v\n", err)
		}
	}

	// Initial start
	startServer()

	// Watch for rebuild signals
	for range rebuild {
		fmt.Printf("\n%s\n", colors.WarningStyle.Render("File changed, reloading..."))
		time.Sleep(100 * time.Millisecond) // Small delay to batch changes
		startServer()
	}
}

func watchFiles(rebuild chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch Go files recursively
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor and .velocity directories
		if strings.Contains(path, "vendor") || strings.Contains(path, ".velocity") {
			return filepath.SkipDir
		}

		// Add directories to watcher
		if info.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Debounce timer
	var debounce *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Only watch .go files
			if !strings.HasSuffix(event.Name, ".go") {
				continue
			}

			// Debounce to avoid multiple rebuilds
			if debounce != nil {
				debounce.Stop()
			}
			debounce = time.AfterFunc(500*time.Millisecond, func() {
				select {
				case rebuild <- true:
				default:
				}
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Watcher error:", err)
		}
	}
}

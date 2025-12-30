package framework

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/ui"
)

// MigrateCmd represents the migrate command
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run all pending database migrations for your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations(false)
	},
}

// MigrateFreshCmd represents the migrate:fresh command
var MigrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Drop all tables and re-run migrations",
	Long:  `Drop all database tables and re-run all migrations from scratch.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations(true)
	},
}

func runMigrations(fresh bool) {
	ui.Header("migrate")

	// Check if we're in a Velocity project
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		ui.Error("Not a Go project (go.mod not found)")
		os.Exit(1)
	}

	// Check if database/migrations directory exists
	if _, err := os.Stat("database/migrations"); os.IsNotExist(err) {
		ui.Error("Migrations directory not found (database/migrations)")
		os.Exit(1)
	}

	// Create migration runner script
	script := `
package main

import (
	"fmt"
	"os"
	"strings"

	_ "IMPORT_PATH/database/migrations"
	"github.com/joho/godotenv"
	"github.com/velocitykode/velocity/pkg/orm"
	"github.com/velocitykode/velocity/pkg/orm/migrate"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorWhite  = "\033[37m"
)

func formatLine(text string) string {
	status := colorGreen + "DONE" + colorReset

	// Calculate dots needed (assuming 180 char width)
	dotsNeeded := 180 - len(text) - 5 // 5 for " DONE"
	if dotsNeeded < 0 {
		dotsNeeded = 3
	}
	dots := strings.Repeat(".", dotsNeeded)

	return fmt.Sprintf("%s %s %s", text, dots, status)
}

func printInfo(message string) {
	fmt.Printf("\n%s%s%s %s\n\n", colorBlue, "[INFO]", colorReset, message)
}

func main() {
	fresh := len(os.Args) > 1 && os.Args[1] == "fresh"

	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}

	if err := orm.InitFromEnv(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	driver := orm.DB()
	if driver == nil {
		fmt.Println("Database driver not initialized")
		os.Exit(1)
	}

	driverName := os.Getenv("DB_CONNECTION")

	migrator := migrate.NewMigrator(driver, driverName)

	fmt.Println(formatLine("Creating migration table"))

	if fresh {
		printInfo("Dropping all tables and re-running migrations.")

		// Get all migrations
		registry := migrate.All()

		if err := migrator.Fresh(); err != nil {
			fmt.Printf("Fresh migration failed: %v\n", err)
			os.Exit(1)
		}

		// Show migrated
		for _, m := range registry {
			text := fmt.Sprintf("%s_%s", m.Version, m.Description)
			fmt.Println(formatLine(text))
		}
	} else {
		// Get pending migrations to show progress
		registry := migrate.All()

		// Get applied migrations
		appliedVersions := make(map[string]bool)
		rows, err := driver.Query("SELECT version FROM migrations")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var version string
				if err := rows.Scan(&version); err == nil {
					appliedVersions[version] = true
				}
			}
		}

		// Find pending
		pending := []migrate.Migration{}
		for _, m := range registry {
			if !appliedVersions[m.Version] {
				pending = append(pending, m)
			}
		}

		if len(pending) == 0 {
			fmt.Println("\nNothing to migrate")
			os.Exit(0)
		}

		printInfo("Running migrations.")

		// Print migrations that will be run
		for _, m := range pending {
			text := fmt.Sprintf("%s_%s", m.Version, m.Description)
			fmt.Println(formatLine(text))
		}

		// Run all migrations via migrator.Up()
		if err := migrator.Up(); err != nil {
			fmt.Printf("\n%sMigration failed:%s %v\n", colorRed, colorReset, err)
			os.Exit(1)
		}
	}

	fmt.Println()
}
`

	// Get module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to read module name: %v", err))
		os.Exit(1)
	}

	// Replace import path
	script = strings.ReplaceAll(script, "IMPORT_PATH", moduleName)

	// Create temporary directory
	tmpDir := ".velocity/tmp"
	os.MkdirAll(tmpDir, 0755)

	// Write temporary migration runner
	tmpFile := fmt.Sprintf("%s/migrate_runner.go", tmpDir)
	if err := os.WriteFile(tmpFile, []byte(script), 0644); err != nil {
		ui.Error(fmt.Sprintf("Failed to create migration runner: %v", err))
		os.Exit(1)
	}
	defer os.Remove(tmpFile)

	// Build
	ui.Step("Compiling migration runner...")
	buildCmd := exec.Command("go", "build", "-o", fmt.Sprintf("%s/migrate", tmpDir), tmpFile)
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		ui.Error(fmt.Sprintf("Build failed: %v\n%s", err, string(buildOutput)))
		os.Exit(1)
	}

	// Run
	var runCmd *exec.Cmd
	if fresh {
		runCmd = exec.Command(fmt.Sprintf("%s/migrate", tmpDir), "fresh")
	} else {
		runCmd = exec.Command(fmt.Sprintf("%s/migrate", tmpDir))
	}
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	if err := runCmd.Run(); err != nil {
		ui.Error(fmt.Sprintf("Migration failed: %v", err))
		os.Exit(1)
	}

	ui.Newline()
	ui.Success("Done")
}

func getModuleName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

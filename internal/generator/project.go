package generator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/velocitykode/velocity-cli/internal/ui"
)

// ProjectConfig holds the configuration for a new project
type ProjectConfig struct {
	Name     string
	Module   string
	Database string
	Cache    string
	Auth     bool
	API      bool
}

// CreateProject generates a new Velocity project from template
func CreateProject(config ProjectConfig) error {
	// Validate project name
	if err := validateProjectName(config.Name); err != nil {
		return err
	}

	// Determine module name
	moduleName := config.Module
	if moduleName == "" {
		moduleName = config.Name
	}

	ui.Info("Creating new Velocity project")
	ui.Step("Cloning project template...")

	// Clone template
	if err := cloneTemplate(config.Name); err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}
	ui.Newline()
	ui.Success("Template cloned")
	ui.Step("Configuring project module...")

	// Replace module name in all files
	if err := replaceModuleName(config.Name, moduleName); err != nil {
		return fmt.Errorf("failed to configure project: %w", err)
	}
	ui.Newline()
	ui.Success("Module configured")
	ui.Step("Initializing Git repository...")

	// Remove template git history and initialize new repo
	if err := reinitGitRepo(config.Name); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}
	ui.Newline()
	ui.Success("Git initialized")
	ui.Step("Creating default migrations...")

	// Create default migrations
	if err := createDefaultMigrations(config.Name); err != nil {
		return fmt.Errorf("failed to create migrations: %w", err)
	}
	ui.Newline()
	ui.Success("Migrations created")
	ui.Step("Setting up environment files...")

	// Create proper .env.example with database config
	if err := createEnvFiles(config); err != nil {
		return fmt.Errorf("failed to create env files: %w", err)
	}
	ui.Newline()
	ui.Success("Environment configured")
	ui.Step("Configuring hot reload...")

	// Setup hot reload
	if err := setupTemplatesAndHotReload(config.Name); err != nil {
		return fmt.Errorf("failed to setup templates: %w", err)
	}
	ui.Newline()
	ui.Success("Hot reload configured")
	ui.Step("Installing dependencies...")

	// Install dependencies
	if err := installDependencies(config.Name); err == nil {
		ui.Newline()
		ui.Success("Dependencies installed")
	}
	ui.Step("Running migrations...")

	// Run migrations
	if err := runMigrations(config.Name); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	ui.Newline()
	ui.Success("Database ready")

	return nil
}

// cloneTemplate clones the velocity-template
func cloneTemplate(projectName string) error {
	// Use git clone directly (gh repo clone can use stale cache)
	cmd := exec.Command("git", "clone", "--depth=1", "git@github.com:velocitykode/velocity-template.git", projectName)
	if err := cmd.Run(); err != nil {
		// Try HTTPS fallback
		cmd = exec.Command("git", "clone", "--depth=1", "https://github.com/velocitykode/velocity-template.git", projectName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone template: %w", err)
		}
	}
	return nil
}

// replaceModuleName replaces {{MODULE_NAME}} in all files
func replaceModuleName(projectPath, moduleName string) error {
	// Get absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("abs path: %w", err)
	}

	// Use find and sed to replace in all Go files
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf("cd '%s' && find . -name '*.go' -type f -exec sed -i '' 's|{{MODULE_NAME}}|%s|g' {} +", absPath, moduleName))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("replace go files: %w: %s", err, string(output))
	}

	// Replace in go.mod
	cmd = exec.Command("sh", "-c",
		fmt.Sprintf("cd '%s' && sed -i '' 's|{{MODULE_NAME}}|%s|g' go.mod", absPath, moduleName))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("replace go.mod: %w: %s", err, string(output))
	}

	// Replace in package.json
	cmd = exec.Command("sh", "-c",
		fmt.Sprintf("cd '%s' && sed -i '' 's|{{MODULE_NAME}}|%s|g' package.json", absPath, moduleName))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("replace package.json: %w: %s", err, string(output))
	}

	// Replace in package-lock.json (if exists)
	cmd = exec.Command("sh", "-c",
		fmt.Sprintf("cd '%s' && [ -f package-lock.json ] && sed -i '' 's|{{MODULE_NAME}}|%s|g' package-lock.json || true", absPath, moduleName))
	cmd.Run() // Ignore error - file may not exist

	// Remove replace directive and fetch pinned version
	cmd = exec.Command("sh", "-c",
		fmt.Sprintf("cd '%s' && sed -i '' '/^replace github.com\\/velocitykode\\/velocity/d' go.mod", absPath))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set pinned version of velocity framework
	cmd = exec.Command("go", "mod", "edit", "-require=github.com/velocitykode/velocity@v0.0.3")
	cmd.Dir = absPath
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set velocity framework version: %w", err)
	}

	return nil
}

// reinitGitRepo removes template git history and creates new repo
func reinitGitRepo(projectPath string) error {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	// Remove .git directory
	gitDir := filepath.Join(absPath, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		return err
	}

	// Initialize new git repo
	originalDir, _ := os.Getwd()
	os.Chdir(absPath)
	defer os.Chdir(originalDir)

	exec.Command("git", "init").Run()
	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()

	return nil
}

// installDependencies runs go mod tidy and bun install
func installDependencies(projectPath string) error {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	originalDir, _ := os.Getwd()
	os.Chdir(absPath)
	defer os.Chdir(originalDir)

	// Run go mod tidy to ensure all dependencies are resolved
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return err
	}

	// Install air for hot reloading
	exec.Command("go", "install", "github.com/air-verse/air@latest").Run()

	// Run bun install (fallback to npm if bun not available)
	if err := exec.Command("bun", "install").Run(); err != nil {
		// Try npm as fallback
		if err := exec.Command("npm", "install").Run(); err != nil {
			return err
		}
	}

	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check for invalid characters
	if strings.ContainsAny(name, " !@#$%^&*()+=[]{}|\\;:'\",<>?/") {
		return fmt.Errorf("project name contains invalid characters")
	}

	// Check if directory already exists
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory %s already exists", name)
	}

	return nil
}

func createDirectoryStructure(projectPath string) error {
	directories := []string{
		"app/http/controllers",
		"app/http/middleware",
		"app/models",
		"bootstrap",
		"config",
		"database/migrations",
		"database/factories",
		"public",
		"resources/views",
		"routes",
		"storage/logs",
		"tests",
	}

	for _, dir := range directories {
		path := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

func initGoModule(config ProjectConfig) error {
	// Change to project directory
	originalDir, _ := os.Getwd()
	os.Chdir(config.Name)
	defer os.Chdir(originalDir)

	// Initialize go module
	moduleName := config.Module
	if moduleName == "" {
		moduleName = config.Name
	}

	cmd := exec.Command("go", "mod", "init", moduleName)
	if err := cmd.Run(); err != nil {
		return err
	}

	ui.Step("Configuring dependencies...")

	// Check if local Velocity exists and use replace directive
	velocityPath := "/Users/ali/code/velocity"
	if _, err := os.Stat(velocityPath); err == nil {
		// Add replace directive for local development
		cmd = exec.Command("go", "mod", "edit", "-replace", "github.com/velocitykode/velocity="+velocityPath)
		cmd.Run()
		ui.Info("Using local Velocity framework")
	} else {
		// Try to get from GitHub (requires GOPRIVATE setup for private repos)
		cmd = exec.Command("go", "get", "github.com/velocitykode/velocity@v0.0.3")
		if err := cmd.Run(); err != nil {
			ui.Warning("Note: Configure GOPRIVATE for private repo access")
		}
	}

	// Add other dependencies based on features
	if config.Database == "postgres" {
		ui.Info("PostgreSQL driver")
		exec.Command("go", "get", "github.com/lib/pq").Run()
	} else if config.Database == "mysql" {
		ui.Info("MySQL driver")
		exec.Command("go", "get", "github.com/go-sql-driver/mysql").Run()
	} else if config.Database == "sqlite" {
		ui.Info("SQLite driver")
		exec.Command("go", "get", "github.com/mattn/go-sqlite3").Run()
	}

	if config.Cache == "redis" {
		ui.Info("Redis client")
		exec.Command("go", "get", "github.com/redis/go-redis/v9").Run()
	}

	// Run go mod tidy
	ui.Step("Tidying up dependencies...")
	exec.Command("go", "mod", "tidy").Run()

	return nil
}

func initGitRepo(projectPath string) {
	originalDir, _ := os.Getwd()
	os.Chdir(projectPath)
	defer os.Chdir(originalDir)

	exec.Command("git", "init").Run()
	exec.Command("git", "add", ".").Run()
}

// InitProject adds Velocity structure to an existing Go project
func InitProject(config ProjectConfig, targetDir string) error {
	ui.Step("Setting up Velocity structure...")
	// Create directory structure in existing directory
	if err := createDirectoryStructure(targetDir); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	ui.Success("Velocity structure created")

	ui.Step("Generating application files...")
	// Generate files from stubs (skip if exists to preserve existing code)
	if err := generateFilesFromStubs(config); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}
	ui.Success("Application files generated")

	ui.Step("Creating configuration files...")
	// Generate config files if they don't exist
	if err := generateProjectFiles(config); err != nil {
		return fmt.Errorf("failed to generate project files: %w", err)
	}
	ui.Success("Configuration files created")

	ui.Step("Adding Velocity dependencies...")
	// Add dependencies to existing go.mod
	if err := addVelocityDependencies(config, targetDir); err != nil {
		return fmt.Errorf("failed to add dependencies: %w", err)
	}
	ui.Success("Dependencies added")

	return nil
}

// addVelocityDependencies adds Velocity and feature dependencies to existing go.mod
// createEnvFiles copies .env.example to .env and generates a new crypto key
func createEnvFiles(config ProjectConfig) error {
	absPath, err := filepath.Abs(config.Name)
	if err != nil {
		return err
	}

	// Copy .env.example to .env
	cmd := exec.Command("cp", filepath.Join(absPath, ".env.example"), filepath.Join(absPath, ".env"))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Generate new crypto key
	newKey, err := generateCryptoKey()
	if err != nil {
		return err
	}

	// Replace crypto key in .env using sed
	cmd = exec.Command("sh", "-c",
		fmt.Sprintf("cd '%s' && sed -i '' 's|^CRYPTO_KEY=.*|CRYPTO_KEY=base64:%s|' .env", absPath, newKey))
	return cmd.Run()
}

// generateCryptoKey generates a new 32-byte base64 encoded key
func generateCryptoKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// createDefaultMigrations creates the 3 default migration files
func createDefaultMigrations(projectPath string) error {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	migrationsDir := filepath.Join(absPath, "database", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return err
	}

	// Migration 1: Create users table
	usersTable := `package migrations

import "github.com/velocitykode/velocity/pkg/orm/migrate"

func init() {
	migrate.Register(&migrate.Migration{
		Version:     "20010101000000",
		Description: "create users table",
		Up: func(m *migrate.Migrator) error {
			return m.CreateTable("users", func(t *migrate.TableBuilder) {
				t.ID()
				t.String("name")
				t.String("email").Unique()
				t.String("password")
				t.String("role").Default("user")
				t.Timestamps()
			})
		},
		Down: func(m *migrate.Migrator) error {
			return m.DropTable("users")
		},
	})
}
`

	// Migration 2: Create cache table
	cacheTable := `package migrations

import "github.com/velocitykode/velocity/pkg/orm/migrate"

func init() {
	migrate.Register(&migrate.Migration{
		Version:     "20010101000001",
		Description: "create cache table",
		Up: func(m *migrate.Migrator) error {
			return m.CreateTable("cache", func(t *migrate.TableBuilder) {
				t.String("key", 255).Unique()
				t.String("value", 10000)
				t.Integer("expiration")
			})
		},
		Down: func(m *migrate.Migrator) error {
			return m.DropTable("cache")
		},
	})
}
`

	// Migration 3: Create jobs table
	jobsTable := `package migrations

import "github.com/velocitykode/velocity/pkg/orm/migrate"

func init() {
	migrate.Register(&migrate.Migration{
		Version:     "20010101000002",
		Description: "create jobs table",
		Up: func(m *migrate.Migrator) error {
			if err := m.CreateTable("jobs", func(t *migrate.TableBuilder) {
				t.ID()
				t.String("queue", 255)
				t.String("payload", 10000)
				t.Integer("attempts").Default("0")
				t.String("scheduled_at", 50)
				t.String("reserved_at", 50).Nullable()
				t.String("reserved_by", 255).Nullable()
				t.String("failed_at", 50).Nullable()
				t.String("failed_reason", 5000).Nullable()
				t.Timestamps()
			}); err != nil {
				return err
			}

			return m.CreateTable("failed_jobs", func(t *migrate.TableBuilder) {
				t.ID()
				t.String("queue", 255)
				t.String("payload", 10000)
				t.String("exception", 10000)
				t.Timestamps()
			})
		},
		Down: func(m *migrate.Migrator) error {
			if err := m.DropTable("failed_jobs"); err != nil {
				return err
			}
			return m.DropTable("jobs")
		},
	})
}
`

	// Write migration files
	migrations := map[string]string{
		"0001_01_01_000000_create_users_table.go": usersTable,
		"0001_01_01_000001_create_cache_table.go": cacheTable,
		"0001_01_01_000002_create_jobs_table.go":  jobsTable,
	}

	for filename, content := range migrations {
		filePath := filepath.Join(migrationsDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// runMigrations runs migrations directly without subprocess
func runMigrations(projectPath string) error {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}

	originalDir, _ := os.Getwd()
	os.Chdir(absPath)
	defer os.Chdir(originalDir)

	// Create temporary migration runner script
	tmpDir := ".velocity/tmp"
	os.MkdirAll(tmpDir, 0755)

	// Get module name
	moduleName, err := getProjectModuleName()
	if err != nil {
		return err
	}

	script := fmt.Sprintf(`
package main

import (
	"fmt"
	"os"

	_ "%s/database/migrations"
	"github.com/joho/godotenv"
	"github.com/velocitykode/velocity/pkg/orm"
	"github.com/velocitykode/velocity/pkg/orm/migrate"
)

const (
	// ANSI symbols (matching CLI style)
	checkSymbol = "\033[32m✓\033[0m"   // green checkmark
	warnSymbol  = "\033[33m!\033[0m"   // yellow warning
	crossSymbol = "\033[31m✗\033[0m"   // red cross
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("%%s .env file not found\n", warnSymbol)
	}

	if err := orm.InitFromEnv(); err != nil {
		fmt.Printf("%%s Failed to initialize database: %%v\n", crossSymbol, err)
		os.Exit(1)
	}

	driver := orm.DB()
	if driver == nil {
		fmt.Printf("%%s Database driver not initialized\n", crossSymbol)
		os.Exit(1)
	}

	driverName := os.Getenv("DB_CONNECTION")
	migrator := migrate.NewMigrator(driver, driverName)

	registry := migrate.All()

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

	pending := []migrate.Migration{}
	for _, m := range registry {
		if !appliedVersions[m.Version] {
			pending = append(pending, m)
		}
	}

	if len(pending) == 0 {
		os.Exit(0)
	}

	if err := migrator.Up(); err != nil {
		fmt.Printf("%%s Migration failed: %%v\n", crossSymbol, err)
		os.Exit(1)
	}

	for _, m := range pending {
		fmt.Printf("%%s \033[32;1m%%s_%%s\033[0m\n", checkSymbol, m.Version, m.Description)
	}
}
`, moduleName)

	tmpFile := fmt.Sprintf("%s/migrate_runner.go", tmpDir)
	if err := os.WriteFile(tmpFile, []byte(script), 0644); err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	// Build
	buildCmd := exec.Command("go", "build", "-o", fmt.Sprintf("%s/migrate", tmpDir), tmpFile)
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build migration runner: %w\n%s", err, string(buildOutput))
	}

	// Run
	runCmd := exec.Command(fmt.Sprintf("%s/migrate", tmpDir))
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	return runCmd.Run()
}

func getProjectModuleName() (string, error) {
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

// StartDevServers starts npm run dev and go run main.go in background
func StartDevServers(projectPath string) {
	const (
		colorReset = "\033[0m"
		colorGreen = "\033[32m"
		colorCyan  = "\033[36m"
	)

	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		fmt.Printf("Failed to resolve project path: %v\n", err)
		return
	}

	// Start npm run dev in background (detached)
	npmCmd := exec.Command("npm", "run", "dev")
	npmCmd.Dir = absPath
	npmCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := npmCmd.Start(); err != nil {
		fmt.Printf("Failed to start npm: %v\n", err)
		return
	}

	// Start air for hot reloading (detached)
	goCmd := exec.Command("air")
	goCmd.Dir = absPath
	goCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := goCmd.Start(); err != nil {
		fmt.Printf("Failed to start Go server: %v\n", err)
		return
	}

	// Show URLs
	ui.KeyValue("Vite", ui.Highlight("http://localhost:5173"))
	ui.KeyValue("Velocity", ui.Highlight("http://localhost:4000"))

	ui.Newline()
	ui.Success("Build something great!")
	ui.Newline()
}

func setupTemplatesAndHotReload(projectPath string) error {
	// .air.toml and tmp/ in .gitignore are now part of the template
	return nil
}

func addVelocityDependencies(config ProjectConfig, projectPath string) error {
	// Change to project directory
	originalDir, _ := os.Getwd()
	os.Chdir(projectPath)
	defer os.Chdir(originalDir)

	ui.Step("Configuring dependencies...")

	// Check if local Velocity exists and use replace directive
	velocityPath := "/Users/ali/code/velocity"
	if _, err := os.Stat(velocityPath); err == nil {
		// Add replace directive for local development
		cmd := exec.Command("go", "mod", "edit", "-replace", "github.com/velocitykode/velocity="+velocityPath)
		cmd.Run()
		ui.Info("Using local Velocity framework")
	} else {
		// Try to get from GitHub
		cmd := exec.Command("go", "get", "github.com/velocitykode/velocity")
		if err := cmd.Run(); err != nil {
			ui.Warning("Note: Configure GOPRIVATE for private repo access")
		}
	}

	// Add other dependencies based on features
	if config.Database == "postgres" {
		ui.Info("PostgreSQL driver")
		exec.Command("go", "get", "github.com/lib/pq").Run()
	} else if config.Database == "mysql" {
		ui.Info("MySQL driver")
		exec.Command("go", "get", "github.com/go-sql-driver/mysql").Run()
	} else if config.Database == "sqlite" {
		ui.Info("SQLite driver")
		exec.Command("go", "get", "github.com/mattn/go-sqlite3").Run()
	}

	if config.Cache == "redis" {
		ui.Info("Redis client")
		exec.Command("go", "get", "github.com/redis/go-redis/v9").Run()
	}

	// Run go mod tidy
	ui.Step("Tidying up dependencies...")
	exec.Command("go", "mod", "tidy").Run()

	return nil
}

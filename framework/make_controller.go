package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/stubs"
	"github.com/velocitykode/velocity-cli/internal/ui"
)

var (
	resource bool
	api      bool
	methods  string
)

// MakeControllerCmd represents the make:controller command
var MakeControllerCmd = &cobra.Command{
	Use:   "make:controller [name]",
	Short: "Generate a new controller",
	Long: `Generate a new HTTP controller with optional CRUD methods.

Examples:
  velocity make:controller UserController
  velocity make:controller PostController --resource
  velocity make:controller API/ProductController --api`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			ui.Error("Controller name is required")
			ui.Newline()
			ui.Muted("Usage: velocity make:controller [name]")
			ui.Newline()
			ui.Muted("Flags:")
			ui.Muted("  --api         Generate API controller (JSON responses)")
			ui.Muted("  --resource    Generate resource controller with CRUD methods")
			ui.Muted("  --methods     Custom methods (comma-separated)")
			return fmt.Errorf("")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		generateController(name)
	},
}

func init() {
	MakeControllerCmd.Flags().BoolVar(&resource, "resource", false, "Generate resource controller with CRUD methods")
	MakeControllerCmd.Flags().BoolVar(&api, "api", false, "Generate API controller (JSON responses)")
	MakeControllerCmd.Flags().StringVar(&methods, "methods", "", "Custom methods (comma-separated)")
}

func generateController(name string) {
	// Parse path and name (no longer require Controller suffix)
	parts := strings.Split(name, "/")
	controllerName := parts[len(parts)-1]

	// Remove Controller suffix if present for cleaner function names
	baseName := strings.TrimSuffix(controllerName, "Controller")

	// Determine file path - using app/http/controllers
	dir := filepath.Join("app", "http", "controllers")
	if len(parts) > 1 {
		subdir := strings.Join(parts[:len(parts)-1], "/")
		dir = filepath.Join(dir, strings.ToLower(subdir))
	}

	// Create directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		ui.Error(fmt.Sprintf("Failed to create directory: %v", err))
		return
	}

	// Generate file path (snake_case with _controller suffix)
	filename := toSnakeCase(baseName) + "_controller.go"
	filePath := filepath.Join(dir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		ui.Error(fmt.Sprintf("Controller already exists: %s", filePath))
		return
	}

	// Create controller file
	file, err := os.Create(filePath)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create file: %v", err))
		return
	}
	defer file.Close()

	// Generate controller code
	data := map[string]interface{}{
		"Package":        getPackageName(dir),
		"ControllerName": baseName,
		"Resource":       resource,
		"API":            api,
		"Methods":        parseCustomMethods(methods),
	}

	tmpl := getControllerTemplate()
	if err := tmpl.Execute(file, data); err != nil {
		ui.Error(fmt.Sprintf("Failed to generate controller: %v", err))
		return
	}

	ui.Success(fmt.Sprintf("Controller created: %s", filePath))
	ui.NextSteps([]string{
		"Register routes in routes/web.go or routes/api.go",
		"Implement controller methods",
	})
}

func getControllerTemplate() *template.Template {
	content, err := stubs.Get("app/http/controllers/controller.go.stub")
	if err != nil {
		// Fallback to basic template if stub not found
		return template.Must(template.New("controller").Parse(`package {{ .Package }}

import "github.com/velocitykode/velocity/pkg/router"

func {{ .ControllerName }}Index(ctx *router.Context) error {
	return ctx.String("Hello from {{ .ControllerName }}")
}
`))
	}
	return template.Must(template.New("controller").Parse(string(content)))
}

func getPackageName(dir string) string {
	parts := strings.Split(dir, string(os.PathSeparator))
	return parts[len(parts)-1]
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func parseCustomMethods(methods string) []string {
	if methods == "" {
		return []string{"Index"}
	}

	parts := strings.Split(methods, ",")
	result := make([]string, 0, len(parts))

	for _, method := range parts {
		method = strings.TrimSpace(method)
		// Capitalize first letter
		if method != "" {
			result = append(result, strings.Title(method))
		}
	}

	return result
}

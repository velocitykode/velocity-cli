package framework

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
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
	// Ensure name ends with Controller
	if !strings.HasSuffix(name, "Controller") {
		name += "Controller"
	}

	// Parse path and name
	parts := strings.Split(name, "/")
	controllerName := parts[len(parts)-1]

	// Determine file path
	dir := filepath.Join("app", "controllers")
	if len(parts) > 1 {
		subdir := strings.Join(parts[:len(parts)-1], "/")
		dir = filepath.Join(dir, strings.ToLower(subdir))
	}

	// Create directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		ui.Error(fmt.Sprintf("Failed to create directory: %v", err))
		return
	}

	// Generate file path
	filename := toSnakeCase(controllerName) + ".go"
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
		"ControllerName": controllerName,
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
	const tmpl = `package {{ .Package }}

import (
	"net/http"
	{{ if .API }}"encoding/json"{{ end }}
)

type {{ .ControllerName }} struct{}

func New{{ .ControllerName }}() *{{ .ControllerName }} {
	return &{{ .ControllerName }}{}
}
{{ if .Resource }}
// Index displays a listing of the resource
func (c *{{ .ControllerName }}) Index(w http.ResponseWriter, r *http.Request) {
	{{ if .API }}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "List all items",
	})
	{{ else }}
	w.Write([]byte("List all items"))
	{{ end }}
}

// Create shows the form for creating a new resource
func (c *{{ .ControllerName }}) Create(w http.ResponseWriter, r *http.Request) {
	{{ if .API }}
	w.WriteHeader(http.StatusMethodNotAllowed)
	{{ else }}
	w.Write([]byte("Show create form"))
	{{ end }}
}

// Store saves a new resource
func (c *{{ .ControllerName }}) Store(w http.ResponseWriter, r *http.Request) {
	{{ if .API }}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Item created",
	})
	{{ else }}
	w.Write([]byte("Store new item"))
	{{ end }}
}

// Show displays the specified resource
func (c *{{ .ControllerName }}) Show(w http.ResponseWriter, r *http.Request) {
	{{ if .API }}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Show item",
	})
	{{ else }}
	w.Write([]byte("Show item"))
	{{ end }}
}

// Edit shows the form for editing the specified resource
func (c *{{ .ControllerName }}) Edit(w http.ResponseWriter, r *http.Request) {
	{{ if .API }}
	w.WriteHeader(http.StatusMethodNotAllowed)
	{{ else }}
	w.Write([]byte("Show edit form"))
	{{ end }}
}

// Update updates the specified resource
func (c *{{ .ControllerName }}) Update(w http.ResponseWriter, r *http.Request) {
	{{ if .API }}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Item updated",
	})
	{{ else }}
	w.Write([]byte("Update item"))
	{{ end }}
}

// Destroy removes the specified resource
func (c *{{ .ControllerName }}) Destroy(w http.ResponseWriter, r *http.Request) {
	{{ if .API }}
	w.WriteHeader(http.StatusNoContent)
	{{ else }}
	w.Write([]byte("Delete item"))
	{{ end }}
}
{{ else }}
{{ range .Methods }}
// {{ . }} handles the {{ . }} action
func (c *{{ $.ControllerName }}) {{ . }}(w http.ResponseWriter, r *http.Request) {
	{{ if $.API }}
	json.NewEncoder(w).Encode(map[string]string{
		"action": "{{ . }}",
	})
	{{ else }}
	w.Write([]byte("{{ . }} action"))
	{{ end }}
}
{{ end }}
{{ end }}
`

	return template.Must(template.New("controller").Parse(tmpl))
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

package framework

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMakeControllerCmd(t *testing.T) {
	// Test command properties
	if MakeControllerCmd.Use != "make:controller [name]" {
		t.Errorf("MakeControllerCmd.Use = %s, want 'make:controller [name]'", MakeControllerCmd.Use)
	}

	if MakeControllerCmd.Short == "" {
		t.Error("MakeControllerCmd.Short is empty")
	}

	// Test that Args is set correctly
	if MakeControllerCmd.Args == nil {
		t.Error("MakeControllerCmd.Args is nil")
	}
}

func TestMakeControllerCmdArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"no args", []string{}, true},
		{"with controller name", []string{"User"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MakeControllerCmd.Args(MakeControllerCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeControllerCmdFlags(t *testing.T) {
	flags := MakeControllerCmd.Flags()

	// Check resource flag
	resourceFlag := flags.Lookup("resource")
	if resourceFlag == nil {
		t.Error("Resource flag not found")
	}

	// Check api flag
	apiFlag := flags.Lookup("api")
	if apiFlag == nil {
		t.Error("API flag not found")
	}

	// Check methods flag
	methodsFlag := flags.Lookup("methods")
	if methodsFlag == nil {
		t.Error("Methods flag not found")
	}
}

func TestGenerateController(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	// Create app/http/controllers directory
	os.MkdirAll("app/http/controllers", 0755)

	// Test generating a basic controller
	generateController("User")

	// Check that controller file was created (snake_case filename)
	controllerPath := filepath.Join("app", "http", "controllers", "user_controller.go")
	if _, err := os.Stat(controllerPath); os.IsNotExist(err) {
		t.Error("Controller file was not created")
	}

	// Read the generated file
	content, err := os.ReadFile(controllerPath)
	if err != nil {
		t.Fatalf("Failed to read controller file: %v", err)
	}

	// Check content
	contentStr := string(content)
	if !strings.Contains(contentStr, "UserIndex") {
		t.Error("Controller does not contain correct function name")
	}

	if !strings.Contains(contentStr, "package controllers") {
		t.Error("Controller does not have correct package")
	}
}

func TestGenerateControllerWithPath(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	// Create app/http/controllers directory
	os.MkdirAll("app/http/controllers", 0755)

	// Test generating a controller in subdirectory
	generateController("API/Product")

	// Check that controller file was created in subdirectory (snake_case filename)
	controllerPath := filepath.Join("app", "http", "controllers", "api", "product_controller.go")
	if _, err := os.Stat(controllerPath); os.IsNotExist(err) {
		t.Error("Controller file was not created in subdirectory")
	}
}

func TestGenerateResourceController(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	// Create app/http/controllers directory
	os.MkdirAll("app/http/controllers", 0755)

	// Set resource flag
	resource = true
	defer func() { resource = false }()

	// Test generating a resource controller
	generateController("Post")

	// Check that controller file was created (snake_case filename)
	controllerPath := filepath.Join("app", "http", "controllers", "post_controller.go")
	content, err := os.ReadFile(controllerPath)
	if err != nil {
		t.Fatalf("Failed to read controller file: %v", err)
	}

	// Check that it contains CRUD methods (function-based naming)
	contentStr := string(content)
	expectedMethods := []string{"PostIndex", "PostCreate", "PostStore", "PostShow", "PostEdit", "PostUpdate", "PostDestroy"}
	for _, method := range expectedMethods {
		if !strings.Contains(contentStr, method) {
			t.Errorf("Resource controller does not contain %s method", method)
		}
	}
}

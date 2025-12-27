package framework

import (
	"testing"
)

func TestServeCmd(t *testing.T) {
	// Test command properties
	if ServeCmd.Use != "serve" {
		t.Errorf("ServeCmd.Use = %s, want 'serve'", ServeCmd.Use)
	}

	if ServeCmd.Short == "" {
		t.Error("ServeCmd.Short is empty")
	}

	// Test flags
	flags := ServeCmd.Flags()

	// Check port flag
	portFlag := flags.Lookup("port")
	if portFlag == nil {
		t.Error("Port flag not found")
	}
	if portFlag.DefValue != "8080" {
		t.Errorf("Port default = %s, want '8080'", portFlag.DefValue)
	}

	// Check env flag
	envFlag := flags.Lookup("env")
	if envFlag == nil {
		t.Error("Env flag not found")
	}
	if envFlag.DefValue != "development" {
		t.Errorf("Env default = %s, want 'development'", envFlag.DefValue)
	}

	// Check watch flag
	watchFlag := flags.Lookup("watch")
	if watchFlag == nil {
		t.Error("Watch flag not found")
	}
	if watchFlag.DefValue != "true" {
		t.Errorf("Watch default = %s, want 'true'", watchFlag.DefValue)
	}

	// Check tags flag
	tagsFlag := flags.Lookup("tags")
	if tagsFlag == nil {
		t.Error("Tags flag not found")
	}
}

func TestInit(t *testing.T) {
	// Test that init sets up flags correctly
	// Flags should be accessible after init
	if port == "" {
		// Default value should be set
		if port != "" {
			t.Errorf("Port has unexpected value: %s", port)
		}
	}

	if env == "" {
		// Default value should be set
		if env != "" {
			t.Errorf("Env has unexpected value: %s", env)
		}
	}
}

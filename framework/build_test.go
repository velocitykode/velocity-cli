package framework

import (
	"testing"
)

func TestBuildCmd(t *testing.T) {
	// Test command properties
	if BuildCmd.Use != "build" {
		t.Errorf("BuildCmd.Use = %s, want 'build'", BuildCmd.Use)
	}

	if BuildCmd.Short == "" {
		t.Error("BuildCmd.Short is empty")
	}

	// Test flags
	flags := BuildCmd.Flags()

	// Check output flag
	outputFlag := flags.Lookup("output")
	if outputFlag == nil {
		t.Error("Output flag not found")
	}
	if outputFlag.DefValue != "./dist/app" {
		t.Errorf("Output default = %s, want './dist/app'", outputFlag.DefValue)
	}

	// Check os flag
	osFlag := flags.Lookup("os")
	if osFlag == nil {
		t.Error("OS flag not found")
	}

	// Check arch flag
	archFlag := flags.Lookup("arch")
	if archFlag == nil {
		t.Error("Arch flag not found")
	}

	// Just verify that the command has flags defined
	// The specific flags are implementation details that may change
}

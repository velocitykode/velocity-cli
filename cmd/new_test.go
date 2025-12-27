package cmd

import (
	"testing"
)

func TestNewCmd(t *testing.T) {
	// Test command properties
	if NewCmd.Use != "new [project-name]" {
		t.Errorf("NewCmd.Use = %s, want 'new [project-name]'", NewCmd.Use)
	}

	if NewCmd.Short == "" {
		t.Error("NewCmd.Short is empty")
	}

	// Test that Args is set correctly
	if NewCmd.Args == nil {
		t.Error("NewCmd.Args is nil")
	}
}

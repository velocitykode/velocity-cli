package framework

import (
	"runtime"
	"testing"
)

func TestBuildCmd(t *testing.T) {
	if BuildCmd.Use != "build" {
		t.Errorf("BuildCmd.Use = %s, want 'build'", BuildCmd.Use)
	}

	if BuildCmd.Short == "" {
		t.Error("BuildCmd.Short is empty")
	}

	// Verify flags exist with correct defaults
	tests := []struct {
		name         string
		defaultValue string
	}{
		{"output", "./dist/app"},
		{"os", runtime.GOOS},
		{"arch", runtime.GOARCH},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := BuildCmd.Flags().Lookup(tt.name)
			if flag == nil {
				t.Errorf("Flag %q not found", tt.name)
				return
			}
			if flag.DefValue != tt.defaultValue {
				t.Errorf("Flag %q default = %q, want %q", tt.name, flag.DefValue, tt.defaultValue)
			}
		})
	}
}

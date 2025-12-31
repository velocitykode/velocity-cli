package framework

import (
	"testing"
)

func TestServeCmd(t *testing.T) {
	if ServeCmd.Use != "serve" {
		t.Errorf("ServeCmd.Use = %s, want 'serve'", ServeCmd.Use)
	}

	if ServeCmd.Short == "" {
		t.Error("ServeCmd.Short is empty")
	}

	// Verify flags exist with correct defaults
	tests := []struct {
		name         string
		defaultValue string
	}{
		{"port", "4000"},
		{"env", "development"},
		{"watch", "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := ServeCmd.Flags().Lookup(tt.name)
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

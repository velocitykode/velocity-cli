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

func TestNewCmdArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"no args", []string{}, true},
		{"with project name", []string{"myproject"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewCmd.Args(NewCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

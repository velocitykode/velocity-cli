package cmd

import (
	"testing"
)

func TestNewCmd_Properties(t *testing.T) {
	if NewCmd.Use != "new [project-name]" {
		t.Errorf("NewCmd.Use = %s, want 'new [project-name]'", NewCmd.Use)
	}

	if NewCmd.Short == "" {
		t.Error("NewCmd.Short is empty")
	}

	if NewCmd.Args == nil {
		t.Error("NewCmd.Args is nil")
	}

	if NewCmd.Run == nil {
		t.Error("NewCmd.Run is nil")
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
		{"with multiple args", []string{"myproject", "extra"}, false},
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

func TestNewCmd_FlagDefaults(t *testing.T) {
	tests := []struct {
		name         string
		defaultValue string
	}{
		{"database", "sqlite"},
		{"cache", "memory"},
		{"auth", "false"},
		{"api", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := NewCmd.Flags().Lookup(tt.name)
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

func TestNewCmd_FlagTypes(t *testing.T) {
	// database and cache are strings
	dbFlag := NewCmd.Flags().Lookup("database")
	if dbFlag.Value.Type() != "string" {
		t.Errorf("database flag type = %s, want string", dbFlag.Value.Type())
	}

	cacheFlag := NewCmd.Flags().Lookup("cache")
	if cacheFlag.Value.Type() != "string" {
		t.Errorf("cache flag type = %s, want string", cacheFlag.Value.Type())
	}

	// auth and api are bools
	authFlag := NewCmd.Flags().Lookup("auth")
	if authFlag.Value.Type() != "bool" {
		t.Errorf("auth flag type = %s, want bool", authFlag.Value.Type())
	}

	apiFlag := NewCmd.Flags().Lookup("api")
	if apiFlag.Value.Type() != "bool" {
		t.Errorf("api flag type = %s, want bool", apiFlag.Value.Type())
	}
}

func TestNewCmd_SilencesErrors(t *testing.T) {
	// NewCmd should silence usage and errors for cleaner output
	if !NewCmd.SilenceUsage {
		t.Error("NewCmd.SilenceUsage should be true")
	}
	if !NewCmd.SilenceErrors {
		t.Error("NewCmd.SilenceErrors should be true")
	}
}

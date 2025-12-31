package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigSetCmdArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"no args", []string{}, true},
		{"one arg", []string{"key"}, true},
		{"two args", []string{"key", "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configSetCmd.Args(configSetCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigGetCmdArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"no args", []string{}, true},
		{"with key", []string{"key"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configGetCmd.Args(configGetCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunConfigSet_ValidKeys(t *testing.T) {
	// Set up temp home dir for config
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	tests := []struct {
		key   string
		value string
	}{
		{"default.database", "postgres"},
		{"default.cache", "redis"},
		{"default.queue", "redis"},
		{"default.auth", "true"},
		{"default.api", "true"},
	}

	cmd := &cobra.Command{}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			err := runConfigSet(cmd, []string{tt.key, tt.value})
			if err != nil {
				t.Errorf("runConfigSet(%q, %q) error = %v", tt.key, tt.value, err)
			}
		})
	}
}

func TestRunConfigSet_InvalidKey(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cmd := &cobra.Command{}
	err := runConfigSet(cmd, []string{"invalid.key", "value"})
	if err == nil {
		t.Error("runConfigSet() should error on invalid key")
	}
	if !strings.Contains(err.Error(), "unknown configuration key") {
		t.Errorf("Error should mention 'unknown configuration key', got: %v", err)
	}
}

func TestRunConfigSet_InvalidDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cmd := &cobra.Command{}
	err := runConfigSet(cmd, []string{"default.database", "invalid_db"})
	if err == nil {
		t.Error("runConfigSet() should error on invalid database value")
	}
}

func TestRunConfigGet_ValidKeys(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cmd := &cobra.Command{}

	// Set a value first
	runConfigSet(cmd, []string{"default.database", "postgres"})

	// Get should not error
	err := runConfigGet(cmd, []string{"default.database"})
	if err != nil {
		t.Errorf("runConfigGet() error = %v", err)
	}
}

func TestRunConfigGet_InvalidKey(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cmd := &cobra.Command{}
	err := runConfigGet(cmd, []string{"invalid.key"})
	if err == nil {
		t.Error("runConfigGet() should error on invalid key")
	}
}

func TestRunConfigList(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cmd := &cobra.Command{}

	// Set some values
	runConfigSet(cmd, []string{"default.database", "postgres"})
	runConfigSet(cmd, []string{"default.cache", "redis"})

	// List should not error
	err := runConfigList(cmd, []string{})
	if err != nil {
		t.Errorf("runConfigList() error = %v", err)
	}
}

func TestRunConfigReset(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cmd := &cobra.Command{}

	// Set a value first (creates config file)
	runConfigSet(cmd, []string{"default.database", "postgres"})

	// Verify config file exists
	configPath := filepath.Join(tmpDir, ".config", "velocity", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("Config file not created, skipping reset test")
	}

	// Reset should not error
	err := runConfigReset(cmd, []string{})
	if err != nil {
		t.Errorf("runConfigReset() error = %v", err)
	}
}

func TestRunConfigReset_NoConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cmd := &cobra.Command{}

	// Reset without existing config should not error
	err := runConfigReset(cmd, []string{})
	if err != nil {
		t.Errorf("runConfigReset() error = %v", err)
	}
}

func TestConfigCmd_HasSubcommands(t *testing.T) {
	subcommands := ConfigCmd.Commands()
	expected := []string{"set", "get", "list", "reset"}

	for _, name := range expected {
		found := false
		for _, cmd := range subcommands {
			if cmd.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ConfigCmd missing subcommand: %s", name)
		}
	}
}

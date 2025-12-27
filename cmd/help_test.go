package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestInitHelp(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
	}

	// Add a test subcommand
	rootCmd.AddCommand(&cobra.Command{
		Use:   "sub",
		Short: "Subcommand",
	})

	InitHelp(rootCmd)

	// Test that help function was set
	if rootCmd.HelpFunc() == nil {
		t.Error("HelpFunc was not set")
	}

	// Test that usage function was set
	if rootCmd.UsageFunc() == nil {
		t.Error("UsageFunc was not set")
	}
}

func TestCustomHelpFunc(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "velocity",
		Short: "Test CLI",
	}

	// Add test commands
	rootCmd.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Test command",
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:    "hidden",
		Short:  "Hidden command",
		Hidden: true,
	})

	// Add flags
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose output")

	InitHelp(rootCmd)

	// Capture output
	oldStdout := rootCmd.OutOrStdout()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	defer rootCmd.SetOut(oldStdout)

	// Execute help
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()

	// Check that output contains expected elements
	expectedElements := []string{
		"██",           // Banner character
		"Usage",        // Usage section
		"Commands",     // Commands section
		"test",         // Test command
		"Test command", // Test command description
		"Flags",        // Flags section
		"verbose",      // Verbose flag
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf("Help output does not contain %q", expected)
		}
	}

	// Check that hidden command is not shown
	if strings.Contains(output, "Hidden command") {
		t.Error("Help output contains hidden command")
	}
}

func TestCustomUsageFunc(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "velocity",
		Short: "Test CLI",
	}

	InitHelp(rootCmd)

	// The usage function should call the help function
	err := rootCmd.UsageFunc()(rootCmd)
	if err != nil {
		t.Errorf("UsageFunc() error = %v", err)
	}
}

package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	// Capture output
	oldStdout := VersionCmd.OutOrStdout()
	buf := new(bytes.Buffer)
	VersionCmd.SetOut(buf)
	defer VersionCmd.SetOut(oldStdout)

	// Execute version command
	VersionCmd.Run(VersionCmd, []string{})

	output := buf.String()

	// Check that output contains version info
	if !strings.Contains(output, "╔") && !strings.Contains(output, "╚") {
		t.Error("Version output does not contain banner box characters")
	}

	if !strings.Contains(output, "0.1.0") {
		t.Errorf("Version output does not contain version number %s", version)
	}
}

func TestVersionValue(t *testing.T) {
	if version == "" {
		t.Error("Version is empty")
	}

	if version != "0.1.0" {
		t.Errorf("Version = %s, want 0.1.0", version)
	}
}

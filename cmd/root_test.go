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

	if !strings.Contains(output, Version) {
		t.Errorf("Version output does not contain version number %s", Version)
	}
}

func TestVersionValue(t *testing.T) {
	if Version == "" {
		t.Error("Version is empty")
	}
}

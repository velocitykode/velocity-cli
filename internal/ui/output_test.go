package ui

import (
	"strings"
	"testing"
)

func TestHighlight(t *testing.T) {
	result := Highlight("test")
	if result == "" {
		t.Error("Highlight should return non-empty string")
	}
	// In non-TTY environments, lipgloss may return unstyled text
	if !strings.Contains(result, "test") {
		t.Error("Highlight should contain the input text")
	}
}

func TestCommand(t *testing.T) {
	result := Command("go run main.go")
	if result == "" {
		t.Error("Command should return non-empty string")
	}
	if !strings.Contains(result, "go run main.go") {
		t.Error("Command should contain the input text")
	}
}

func TestStylesExist(t *testing.T) {
	// Test that styling functions return non-empty strings
	if Highlight("test") == "" {
		t.Error("Highlight should return non-empty string")
	}
	if Command("test") == "" {
		t.Error("Command should return non-empty string")
	}
}

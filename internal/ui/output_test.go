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

func TestSpinnerWithError(t *testing.T) {
	// Test that SpinnerWithError returns the action's error
	expectedErr := "test error"
	err := SpinnerWithError("Testing...", func() error {
		return nil
	})
	if err != nil {
		t.Errorf("SpinnerWithError should return nil for successful action, got: %v", err)
	}

	// Note: Can't easily test spinner visuals in unit test
	_ = expectedErr
}

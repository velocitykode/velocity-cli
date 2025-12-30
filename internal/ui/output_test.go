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

func TestLoader(t *testing.T) {
	called := false
	err := Loader("Testing...", func() error {
		called = true
		return nil
	})
	if err != nil {
		t.Errorf("Loader should return nil for successful action, got: %v", err)
	}
	if !called {
		t.Error("Loader action was not called")
	}
}

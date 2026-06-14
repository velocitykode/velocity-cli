package prism

import (
	"strings"
	"testing"
)

func TestNote(t *testing.T) {
	out := capture(func() { Note("Server running on port 4000") })
	if !strings.Contains(out, "Server running on port 4000") {
		t.Errorf("Note() output missing message, got: %q", out)
	}
}

func TestAlert(t *testing.T) {
	out := capture(func() { Alert("Connection failed") })
	if !strings.Contains(out, "Connection failed") {
		t.Errorf("Alert() output missing message, got: %q", out)
	}
}

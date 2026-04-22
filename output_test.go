package cli

import (
	"bytes"
	"strings"
	"testing"
)

func capture(fn func()) string {
	var buf bytes.Buffer
	SetWriter(&buf)
	defer SetWriter(nil)
	fn()
	return buf.String()
}

func TestInfo(t *testing.T) {
	out := capture(func() { Info("hello world") })
	if !strings.Contains(out, "hello world") {
		t.Errorf("Info() output missing message, got: %q", out)
	}
}

func TestTip(t *testing.T) {
	out := capture(func() { Tip("install bun") })
	if !strings.Contains(out, "install bun") {
		t.Errorf("Tip() output missing message, got: %q", out)
	}
	if !strings.Contains(out, "ℹ") {
		t.Errorf("Tip() output missing info glyph, got: %q", out)
	}
}

func TestSuccess(t *testing.T) {
	out := capture(func() { Success("done") })
	if !strings.Contains(out, "done") {
		t.Errorf("Success() output missing message, got: %q", out)
	}
}

func TestWarning(t *testing.T) {
	out := capture(func() { Warning("careful") })
	if !strings.Contains(out, "careful") {
		t.Errorf("Warning() output missing message, got: %q", out)
	}
}

func TestError(t *testing.T) {
	out := capture(func() { Error("failed") })
	if !strings.Contains(out, "failed") {
		t.Errorf("Error() output missing message, got: %q", out)
	}
}

func TestMuted(t *testing.T) {
	out := capture(func() { Muted("dim text") })
	if !strings.Contains(out, "dim text") {
		t.Errorf("Muted() output missing message, got: %q", out)
	}
}

func TestBold(t *testing.T) {
	out := capture(func() { Bold("strong") })
	if !strings.Contains(out, "strong") {
		t.Errorf("Bold() output missing message, got: %q", out)
	}
}

func TestHeader(t *testing.T) {
	out := capture(func() { Header("migrate") })
	if !strings.Contains(out, "MIGRATE") {
		t.Errorf("Header() should uppercase, got: %q", out)
	}
}

func TestNewline(t *testing.T) {
	out := capture(func() { Newline() })
	if out != "\n" {
		t.Errorf("Newline() expected newline, got: %q", out)
	}
}

func TestHighlight(t *testing.T) {
	result := Highlight("4000")
	if !strings.Contains(result, "4000") {
		t.Errorf("Highlight() should contain text, got: %q", result)
	}
}

func TestKeyValue(t *testing.T) {
	out := capture(func() { KeyValue("port", "4000") })
	if !strings.Contains(out, "port:") || !strings.Contains(out, "4000") {
		t.Errorf("KeyValue() output missing key or value, got: %q", out)
	}
}

func TestStep(t *testing.T) {
	out := capture(func() { Step("building...") })
	if !strings.Contains(out, "building...") {
		t.Errorf("Step() output missing message, got: %q", out)
	}
}

func TestNextSteps(t *testing.T) {
	out := capture(func() { NextSteps([]string{"run migrate", "start server"}) })
	if !strings.Contains(out, "run migrate") || !strings.Contains(out, "start server") {
		t.Errorf("NextSteps() output missing steps, got: %q", out)
	}
}

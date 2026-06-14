package prism

import (
	"bytes"
	"os"
	"testing"
)

func TestSetWriter_Override(t *testing.T) {
	var buf bytes.Buffer
	SetWriter(&buf)
	defer SetWriter(nil)
	fprintln("test")
	if buf.String() != "test\n" {
		t.Errorf("expected 'test\\n', got %q", buf.String())
	}
}

func TestSetWriter_NilResolvesStdoutLazily(t *testing.T) {
	// Default state: output is nil, writer() should return current os.Stdout.
	SetWriter(nil)
	if w := writer(); w != os.Stdout {
		t.Errorf("expected writer() to return os.Stdout, got %T", w)
	}

	// Swap os.Stdout mid-test; writer() should pick up the swap.
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w
	if got := writer(); got != w {
		t.Error("writer() should resolve to the swapped os.Stdout")
	}
	_ = r
	w.Close()
}

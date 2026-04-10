package cli

import (
	"bytes"
	"os"
	"testing"
)

func TestSetWriter(t *testing.T) {
	var buf bytes.Buffer
	SetWriter(&buf)
	fprintln("test")
	if buf.String() != "test\n" {
		t.Errorf("expected 'test\\n', got %q", buf.String())
	}

	SetWriter(nil)
	if output != os.Stdout {
		t.Error("SetWriter(nil) should reset to os.Stdout")
	}
}

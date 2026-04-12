package cli

import (
	"fmt"
	"io"
	"os"
)

// output is the override writer set by SetWriter. When nil (the default),
// output resolves to os.Stdout at call time - so consumers swapping
// os.Stdout in tests work transparently, matching the fmt package.
var output io.Writer

// SetWriter sets the output writer for all non-interactive components.
// Pass nil to fall back to os.Stdout (resolved dynamically at each call).
func SetWriter(w io.Writer) {
	output = w
}

func writer() io.Writer {
	if output == nil {
		return os.Stdout
	}
	return output
}

func fprintln(a ...interface{}) {
	fmt.Fprintln(writer(), a...)
}

func fprintf(format string, a ...interface{}) {
	fmt.Fprintf(writer(), format, a...)
}

package cli

import (
	"fmt"
	"io"
	"os"
)

var output io.Writer = os.Stdout

// SetWriter sets the output writer for all non-interactive components.
// Pass nil to reset to os.Stdout.
func SetWriter(w io.Writer) {
	if w == nil {
		output = os.Stdout
		return
	}
	output = w
}

func fprintln(a ...interface{}) {
	fmt.Fprintln(output, a...)
}

func fprintf(format string, a ...interface{}) {
	fmt.Fprintf(output, format, a...)
}

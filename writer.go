package prism

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

// interactive reports whether the process is attached to a terminal capable of
// driving the bubbletea-based components. Both stdin (bubbletea reads key
// input from the controlling TTY) and stderr (where animations render) must be
// terminals; otherwise components must degrade to plain, non-TTY output rather
// than failing to open /dev/tty in CI/headless/piped contexts.
func interactive() bool {
	return isTerminal(os.Stdin.Fd()) && isTerminal(os.Stderr.Fd())
}

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

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

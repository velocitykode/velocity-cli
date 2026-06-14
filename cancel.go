package prism

import "os"

// Canceller is implemented by interactive prompt models that can be
// cancelled via Ctrl+C or Esc. Implementations return true once the
// user has signalled a cancel.
type Canceller interface {
	Cancelled() bool
}

var cancelHandler func()

// SetCancelHandler registers a callback invoked right before the
// process exits on cancel. Pass nil to restore the default (exit only).
// The callback should be cheap and non-blocking; it runs inside the
// prompt's return path.
func SetCancelHandler(fn func()) {
	cancelHandler = fn
}

// ExitOnCancel inspects the final prompt model and, if the user
// cancelled, runs the registered handler (if any) and exits the
// process with code 1. Matches laravel/prompts semantics: exit(1) by
// default, overridable via SetCancelHandler.
//
// Prompt implementations should call this immediately after the
// bubbletea program returns, before reading the model's value.
var exitFn = os.Exit

func ExitOnCancel(m any) {
	c, ok := m.(Canceller)
	if !ok || !c.Cancelled() {
		return
	}
	if cancelHandler != nil {
		cancelHandler()
	}
	exitFn(1)
}

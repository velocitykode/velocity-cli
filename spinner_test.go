package cli

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSpinnerModel_Init(t *testing.T) {
	m := newSpinnerModel("loading")
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() should return a tick command")
	}
}

func TestSpinnerModel_DoneMsg(t *testing.T) {
	m := newSpinnerModel("loading")
	updated, cmd := m.Update(spinnerDoneMsg{err: nil})
	sm := updated.(spinnerModel)
	if !sm.done {
		t.Error("expected done=true after spinnerDoneMsg")
	}
	if sm.err != nil {
		t.Error("expected nil error")
	}
	if cmd == nil {
		t.Error("expected tea.Quit command")
	}
}

func TestSpinnerModel_DoneWithError(t *testing.T) {
	m := newSpinnerModel("loading")
	e := errors.New("fail")
	updated, _ := m.Update(spinnerDoneMsg{err: e})
	sm := updated.(spinnerModel)
	if sm.err != e {
		t.Errorf("expected error %v, got %v", e, sm.err)
	}
}

func TestSpinnerModel_CtrlC(t *testing.T) {
	m := newSpinnerModel("loading")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("ctrl+c should quit")
	}
}

// TestSpinner_NonInteractive_RunsFnAndReturnsError covers the headless/CI
// path: with no terminal (as under `go test`), Spinner must run fn and return
// its result verbatim, never the bubbletea "could not open a new TTY" error.
func TestSpinner_NonInteractive_RunsFnAndReturnsError(t *testing.T) {
	called := false
	if err := Spinner("working", func() error { called = true; return nil }); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !called {
		t.Fatal("fn was not executed")
	}

	want := errors.New("boom")
	if got := Spinner("working", func() error { return want }); got != want {
		t.Fatalf("expected fn error %v, got %v", want, got)
	}
}

func TestSpinnerModel_View(t *testing.T) {
	m := newSpinnerModel("building")
	v := m.View()
	if v == "" {
		t.Error("View() should not be empty when not done")
	}

	m.done = true
	v = m.View()
	if v != "" {
		t.Error("View() should be empty when done")
	}
}

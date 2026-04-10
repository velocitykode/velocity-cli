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

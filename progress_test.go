package cli

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestProgressModel_Init(t *testing.T) {
	m := newProgressModel(10)
	cmd := m.Init()
	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestProgressModel_Increment(t *testing.T) {
	m := newProgressModel(3)
	updated, _ := m.Update(progressIncrMsg{})
	pm := updated.(progressModel)
	if pm.completed != 1 {
		t.Errorf("expected completed=1, got %d", pm.completed)
	}
}

func TestProgressModel_Complete(t *testing.T) {
	m := newProgressModel(1)
	updated, cmd := m.Update(progressIncrMsg{})
	pm := updated.(progressModel)
	if pm.completed != 1 {
		t.Errorf("expected completed=1, got %d", pm.completed)
	}
	if cmd == nil {
		t.Error("expected quit command on completion")
	}
}

func TestProgressModel_KeyQuit(t *testing.T) {
	m := newProgressModel(5)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("expected quit on key press")
	}
}

func TestProgressModel_View(t *testing.T) {
	m := newProgressModel(10)
	v := m.View()
	if v == "" {
		t.Error("View() should not be empty")
	}
}

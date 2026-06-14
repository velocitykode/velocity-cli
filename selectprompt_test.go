package prism

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSelectModel_Navigation(t *testing.T) {
	m := newSelectModel("Pick:", []string{"a", "b", "c"}, &selectConfig{})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm := updated.(selectModel)
	if sm.cursor != 1 {
		t.Errorf("down: expected cursor=1, got %d", sm.cursor)
	}

	updated, _ = sm.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm = updated.(selectModel)
	if sm.cursor != 2 {
		t.Errorf("down again: expected cursor=2, got %d", sm.cursor)
	}

	// Should not go past last
	updated, _ = sm.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm = updated.(selectModel)
	if sm.cursor != 2 {
		t.Errorf("down at end: expected cursor=2, got %d", sm.cursor)
	}

	updated, _ = sm.Update(tea.KeyMsg{Type: tea.KeyUp})
	sm = updated.(selectModel)
	if sm.cursor != 1 {
		t.Errorf("up: expected cursor=1, got %d", sm.cursor)
	}
}

func TestSelectModel_Enter(t *testing.T) {
	m := newSelectModel("Pick:", []string{"a", "b", "c"}, &selectConfig{})
	m.cursor = 1
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	sm := updated.(selectModel)
	if !sm.done {
		t.Error("enter should set done")
	}
	if sm.value() != "b" {
		t.Errorf("expected 'b', got %q", sm.value())
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestSelectModel_Default(t *testing.T) {
	m := newSelectModel("DB:", []string{"mysql", "postgres", "sqlite"}, &selectConfig{defaultVal: "postgres"})
	if m.cursor != 1 {
		t.Errorf("expected cursor at postgres (1), got %d", m.cursor)
	}
}

func TestSelectModel_VimKeys(t *testing.T) {
	m := newSelectModel("Pick:", []string{"a", "b", "c"}, &selectConfig{})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	sm := updated.(selectModel)
	if sm.cursor != 1 {
		t.Error("j should move down")
	}

	updated, _ = sm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	sm = updated.(selectModel)
	if sm.cursor != 0 {
		t.Error("k should move up")
	}
}

func TestSelectModel_View(t *testing.T) {
	m := newSelectModel("Pick:", []string{"a", "b"}, &selectConfig{})
	v := m.View()
	if !strings.Contains(v, "Pick:") {
		t.Error("view should show label")
	}
	if !strings.Contains(v, "a") || !strings.Contains(v, "b") {
		t.Error("view should show options")
	}
}

func TestSelectModel_CtrlC(t *testing.T) {
	m := newSelectModel("Pick:", []string{"a", "b"}, &selectConfig{})
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	sm := updated.(selectModel)
	if !sm.Cancelled() {
		t.Error("ctrl+c should mark model cancelled")
	}
	if !sm.done {
		t.Error("ctrl+c should set done")
	}
	if sm.value() != "" {
		t.Error("ctrl+c should return empty value")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

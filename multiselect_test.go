package cli

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMultiselectModel_Toggle(t *testing.T) {
	m := newMultiselectModel("Pick:", []string{"a", "b", "c"})

	// Toggle first
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	mm := updated.(multiselectModel)
	if !mm.selected[0] {
		t.Error("space should toggle selection on")
	}

	// Toggle off
	updated, _ = mm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	mm = updated.(multiselectModel)
	if mm.selected[0] {
		t.Error("space again should toggle selection off")
	}
}

func TestMultiselectModel_Navigation(t *testing.T) {
	m := newMultiselectModel("Pick:", []string{"a", "b", "c"})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	mm := updated.(multiselectModel)
	if mm.cursor != 1 {
		t.Errorf("expected cursor=1, got %d", mm.cursor)
	}
}

func TestMultiselectModel_ToggleAll(t *testing.T) {
	m := newMultiselectModel("Pick:", []string{"a", "b", "c"})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	mm := updated.(multiselectModel)
	vals := mm.values()
	if len(vals) != 3 {
		t.Errorf("toggle all: expected 3 selected, got %d", len(vals))
	}

	// Toggle all off
	updated, _ = mm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	mm = updated.(multiselectModel)
	vals = mm.values()
	if len(vals) != 0 {
		t.Errorf("toggle all off: expected 0 selected, got %d", len(vals))
	}
}

func TestMultiselectModel_Enter(t *testing.T) {
	m := newMultiselectModel("Pick:", []string{"a", "b", "c"})
	m.selected[0] = true
	m.selected[2] = true

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	mm := updated.(multiselectModel)
	if !mm.done {
		t.Error("enter should set done")
	}
	vals := mm.values()
	if len(vals) != 2 || vals[0] != "a" || vals[1] != "c" {
		t.Errorf("expected [a, c], got %v", vals)
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestMultiselectModel_CtrlC(t *testing.T) {
	m := newMultiselectModel("Pick:", []string{"a", "b"})
	m.selected[0] = true

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	mm := updated.(multiselectModel)
	if !mm.Cancelled() {
		t.Error("ctrl+c should mark model cancelled")
	}
	if mm.selected != nil {
		t.Error("ctrl+c should clear selection")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

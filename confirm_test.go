package cli

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestConfirmModel_DefaultNo(t *testing.T) {
	m := newConfirmModel("Delete?", &confirmConfig{defaultVal: false})
	if m.value != false {
		t.Error("default should be false")
	}
}

func TestConfirmModel_DefaultYes(t *testing.T) {
	m := newConfirmModel("Continue?", &confirmConfig{defaultVal: true})
	if m.value != true {
		t.Error("default should be true")
	}
}

func TestConfirmModel_PressY(t *testing.T) {
	m := newConfirmModel("Delete?", &confirmConfig{})
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	cm := updated.(confirmModel)
	if !cm.value || !cm.done {
		t.Error("pressing 'y' should set value=true and done=true")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestConfirmModel_PressN(t *testing.T) {
	m := newConfirmModel("Delete?", &confirmConfig{defaultVal: true})
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	cm := updated.(confirmModel)
	if cm.value || !cm.done {
		t.Error("pressing 'n' should set value=false and done=true")
	}
}

func TestConfirmModel_EnterUsesDefault(t *testing.T) {
	m := newConfirmModel("Continue?", &confirmConfig{defaultVal: true})
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	cm := updated.(confirmModel)
	if !cm.value {
		t.Error("enter should use default value (true)")
	}
}

func TestConfirmModel_Toggle(t *testing.T) {
	m := newConfirmModel("Delete?", &confirmConfig{defaultVal: false})
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	cm := updated.(confirmModel)
	if cm.value != true {
		t.Error("right arrow should toggle to true")
	}
}

func TestConfirmModel_ViewShowsHint(t *testing.T) {
	m := newConfirmModel("Delete?", &confirmConfig{defaultVal: true})
	v := m.View()
	if !strings.Contains(v, "Y/n") {
		t.Errorf("view should show [Y/n] for defaultYes, got: %q", v)
	}
}

func TestConfirmModel_ViewDone(t *testing.T) {
	m := newConfirmModel("Delete?", &confirmConfig{})
	m.done = true
	m.value = true
	v := m.View()
	if !strings.Contains(v, "Yes") {
		t.Errorf("done view should show 'Yes', got: %q", v)
	}
}

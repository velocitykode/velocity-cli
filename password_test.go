package cli

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPasswordModel_MaskedView(t *testing.T) {
	m := newPasswordModel("Password:")
	m.textInput.SetValue("secret")
	v := m.View()
	if strings.Contains(v, "secret") {
		t.Error("password should not show plain text")
	}
	if !strings.Contains(v, "Password:") {
		t.Error("view should show label")
	}
}

func TestPasswordModel_Enter(t *testing.T) {
	m := newPasswordModel("Password:")
	m.textInput.SetValue("mypass")
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	pm := updated.(passwordModel)
	if !pm.done {
		t.Error("enter should set done")
	}
	if pm.textInput.Value() != "mypass" {
		t.Errorf("expected 'mypass', got %q", pm.textInput.Value())
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestPasswordModel_DoneView(t *testing.T) {
	m := newPasswordModel("Password:")
	m.done = true
	v := m.View()
	if !strings.Contains(v, "••••••••") {
		t.Errorf("done view should show dots, got: %q", v)
	}
}

package prism

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTextModel_Init(t *testing.T) {
	m := newTextModel("Name:", &textConfig{})
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() should return blink command")
	}
}

func TestTextModel_RequiredEmpty(t *testing.T) {
	m := newTextModel("Name:", &textConfig{required: true})
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	tm := updated.(textModel)
	if tm.done {
		t.Error("should not be done when required field is empty")
	}
	if tm.errMsg == "" {
		t.Error("should have error message")
	}
}

func TestTextModel_ValidationError(t *testing.T) {
	cfg := &textConfig{
		validate: func(s string) error {
			if len(s) < 3 {
				return errors.New("too short")
			}
			return nil
		},
	}
	m := newTextModel("Name:", cfg)
	m.textInput.SetValue("ab")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	tm := updated.(textModel)
	if tm.done {
		t.Error("should not be done with validation error")
	}
	if tm.errMsg != "too short" {
		t.Errorf("expected 'too short', got %q", tm.errMsg)
	}
}

func TestTextModel_ValidationPass(t *testing.T) {
	cfg := &textConfig{
		validate: func(s string) error {
			if len(s) < 3 {
				return errors.New("too short")
			}
			return nil
		},
	}
	m := newTextModel("Name:", cfg)
	m.textInput.SetValue("alice")
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	tm := updated.(textModel)
	if !tm.done {
		t.Error("should be done when validation passes")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestTextModel_DefaultValue(t *testing.T) {
	m := newTextModel("Name:", &textConfig{defaultVal: "World"})
	if m.textInput.Value() != "World" {
		t.Errorf("expected default 'World', got %q", m.textInput.Value())
	}
}

func TestTextModel_Value(t *testing.T) {
	m := newTextModel("Name:", &textConfig{defaultVal: "fallback"})
	m.textInput.SetValue("alice")
	if m.value() != "alice" {
		t.Errorf("expected 'alice', got %q", m.value())
	}

	m.textInput.SetValue("")
	if m.value() != "fallback" {
		t.Errorf("expected fallback 'fallback', got %q", m.value())
	}
}

func TestTextModel_ViewShowsLabel(t *testing.T) {
	m := newTextModel("Model name:", &textConfig{})
	v := m.View()
	if !strings.Contains(v, "Model name:") {
		t.Errorf("view should show label, got: %q", v)
	}
}

func TestTextModel_CtrlC(t *testing.T) {
	m := newTextModel("Name:", &textConfig{})
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm := updated.(textModel)
	if !tm.done {
		t.Error("ctrl+c should set done")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TextOption configures a Text prompt.
type TextOption func(*textConfig)

type textConfig struct {
	placeholder string
	defaultVal  string
	required    bool
	validate    func(string) error
}

// WithPlaceholder sets placeholder text shown when input is empty.
func WithPlaceholder(s string) TextOption {
	return func(c *textConfig) { c.placeholder = s }
}

// WithDefault sets a default value pre-filled in the input.
func WithDefault(s string) TextOption {
	return func(c *textConfig) { c.defaultVal = s }
}

// WithRequired marks the input as required.
func WithRequired() TextOption {
	return func(c *textConfig) { c.required = true }
}

// WithValidation sets a validation function.
func WithValidation(fn func(string) error) TextOption {
	return func(c *textConfig) { c.validate = fn }
}

type textModel struct {
	label     string
	textInput textinput.Model
	cfg       *textConfig
	errMsg    string
	done      bool
}

func newTextModel(label string, cfg *textConfig) textModel {
	ti := textinput.New()
	ti.Placeholder = cfg.placeholder
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40
	ti.PromptStyle = cursorStyle
	ti.TextStyle = boldStyle

	if cfg.defaultVal != "" {
		ti.SetValue(cfg.defaultVal)
	}

	return textModel{
		label:     label,
		textInput: ti,
		cfg:       cfg,
	}
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.done = true
			return m, tea.Quit
		case tea.KeyEnter:
			val := m.textInput.Value()
			if m.cfg.required && val == "" {
				m.errMsg = "This field is required."
				return m, nil
			}
			if m.cfg.validate != nil {
				if err := m.cfg.validate(val); err != nil {
					m.errMsg = err.Error()
					return m, nil
				}
			}
			m.errMsg = ""
			m.done = true
			return m, tea.Quit
		}
	}

	m.errMsg = ""
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m textModel) View() string {
	if m.done {
		val := m.textInput.Value()
		if val == "" {
			val = m.cfg.defaultVal
		}
		return fmt.Sprintf("  %s %s\n", promptLabelStyle.Render(m.label), successStyle.Render(val)) + "\n"
	}

	s := fmt.Sprintf("  %s\n  %s\n", promptLabelStyle.Render(m.label), m.textInput.View())
	if m.errMsg != "" {
		s += fmt.Sprintf("  %s\n", errorMsgStyle.Render(m.errMsg))
	}
	s += fmt.Sprintf("\n  %s\n", promptHelpStyle.Render("enter to confirm"))
	return s
}

func (m textModel) value() string {
	v := m.textInput.Value()
	if v == "" {
		return m.cfg.defaultVal
	}
	return v
}

// Text displays a single-line text input and returns the entered value.
func Text(label string, opts ...TextOption) string {
	cfg := &textConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	m := newTextModel(label, cfg)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return cfg.defaultVal
	}
	if fm, ok := finalModel.(textModel); ok {
		return fm.value()
	}
	return cfg.defaultVal
}

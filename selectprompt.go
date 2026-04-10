package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// SelectOption configures a Select prompt.
type SelectOption func(*selectConfig)

type selectConfig struct {
	defaultVal string
}

// WithSelectDefault sets the initially highlighted option.
func WithSelectDefault(val string) SelectOption {
	return func(c *selectConfig) { c.defaultVal = val }
}

type selectModel struct {
	label   string
	options []string
	cursor  int
	done    bool
}

func newSelectModel(label string, options []string, cfg *selectConfig) selectModel {
	cursor := 0
	if cfg.defaultVal != "" {
		for i, o := range options {
			if o == cfg.defaultVal {
				cursor = i
				break
			}
		}
	}
	return selectModel{
		label:   label,
		options: options,
		cursor:  cursor,
	}
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cursor = -1
			m.done = true
			return m, tea.Quit
		case tea.KeyEnter:
			m.done = true
			return m, tea.Quit
		case tea.KeyUp, tea.KeyShiftTab:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown, tea.KeyTab:
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		}

		switch msg.String() {
		case "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	if m.done {
		val := ""
		if m.cursor >= 0 && m.cursor < len(m.options) {
			val = m.options[m.cursor]
		}
		return fmt.Sprintf("  %s %s\n", promptLabelStyle.Render(m.label), successStyle.Render(val)) + "\n"
	}

	s := fmt.Sprintf("  %s\n\n", promptLabelStyle.Render(m.label))
	for i, opt := range m.options {
		if i == m.cursor {
			s += fmt.Sprintf("  %s %s\n", cursorStyle.Render("▸"), selectedStyle.Render(opt))
		} else {
			s += fmt.Sprintf("    %s\n", unselectedStyle.Render(opt))
		}
	}
	s += fmt.Sprintf("\n  %s\n", promptHelpStyle.Render("↑/↓ to navigate, enter to select"))
	return s
}

func (m selectModel) value() string {
	if m.cursor >= 0 && m.cursor < len(m.options) {
		return m.options[m.cursor]
	}
	return ""
}

// Select displays a list of options and returns the selected one.
func Select(label string, options []string, opts ...SelectOption) string {
	cfg := &selectConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	m := newSelectModel(label, options, cfg)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return ""
	}
	if fm, ok := finalModel.(selectModel); ok {
		return fm.value()
	}
	return ""
}

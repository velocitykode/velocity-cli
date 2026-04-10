package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmOption configures a Confirm prompt.
type ConfirmOption func(*confirmConfig)

type confirmConfig struct {
	defaultVal bool
}

// WithDefaultYes sets the default to yes.
func WithDefaultYes() ConfirmOption {
	return func(c *confirmConfig) { c.defaultVal = true }
}

// WithDefaultNo sets the default to no.
func WithDefaultNo() ConfirmOption {
	return func(c *confirmConfig) { c.defaultVal = false }
}

type confirmModel struct {
	label      string
	defaultVal bool
	value      bool
	done       bool
}

func newConfirmModel(label string, cfg *confirmConfig) confirmModel {
	return confirmModel{
		label:      label,
		defaultVal: cfg.defaultVal,
		value:      cfg.defaultVal,
	}
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.value = false
			m.done = true
			return m, tea.Quit
		case tea.KeyEnter:
			m.done = true
			return m, tea.Quit
		case tea.KeyLeft, tea.KeyRight:
			m.value = !m.value
			return m, nil
		}

		switch strings.ToLower(msg.String()) {
		case "y":
			m.value = true
			m.done = true
			return m, tea.Quit
		case "n":
			m.value = false
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	if m.done {
		choice := "No"
		if m.value {
			choice = "Yes"
		}
		return fmt.Sprintf("  %s %s\n", promptLabelStyle.Render(m.label), successStyle.Render(choice)) + "\n"
	}

	hint := "[y/N]"
	if m.defaultVal {
		hint = "[Y/n]"
	}

	yes := unselectedStyle.Render("Yes")
	no := unselectedStyle.Render("No")
	if m.value {
		yes = selectedStyle.Render("▸ Yes")
		no = unselectedStyle.Render("  No")
	} else {
		yes = unselectedStyle.Render("  Yes")
		no = selectedStyle.Render("▸ No")
	}

	return fmt.Sprintf(
		"  %s %s\n\n  %s  %s\n\n  %s\n",
		promptLabelStyle.Render(m.label),
		promptHelpStyle.Render(hint),
		yes, no,
		promptHelpStyle.Render("← / → to toggle, enter to confirm"),
	)
}

// Confirm displays a yes/no prompt and returns the user's choice.
func Confirm(label string, opts ...ConfirmOption) bool {
	cfg := &confirmConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	m := newConfirmModel(label, cfg)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return cfg.defaultVal
	}
	if fm, ok := finalModel.(confirmModel); ok {
		return fm.value
	}
	return cfg.defaultVal
}

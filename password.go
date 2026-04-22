package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type passwordModel struct {
	label     string
	textInput textinput.Model
	done      bool
	cancelled bool
}

func (m passwordModel) Cancelled() bool { return m.cancelled }

func newPasswordModel(label string) passwordModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.PromptStyle = cursorStyle

	return passwordModel{
		label:     label,
		textInput: ti,
	}
}

func (m passwordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m passwordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.cancelled = true
			m.done = true
			return m, tea.Quit
		case tea.KeyEnter:
			m.done = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m passwordModel) View() string {
	if m.done {
		return fmt.Sprintf("  %s %s\n", promptLabelStyle.Render(m.label), mutedStyle.Render("••••••••")) + "\n"
	}

	return fmt.Sprintf(
		"  %s\n  %s\n\n  %s\n",
		promptLabelStyle.Render(m.label),
		m.textInput.View(),
		promptHelpStyle.Render("enter to confirm"),
	)
}

// Password displays a masked text input and returns the entered value.
func Password(label string) string {
	m := newPasswordModel(label)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return ""
	}
	ExitOnCancel(finalModel)
	if fm, ok := finalModel.(passwordModel); ok {
		return fm.textInput.Value()
	}
	return ""
}

package cli

import (
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type spinnerDoneMsg struct{ err error }

type spinnerModel struct {
	spinner spinner.Model
	message string
	err     error
	done    bool
}

func newSpinnerModel(message string) spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorPrimary)
	return spinnerModel{spinner: s, message: message}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case spinnerDoneMsg:
		m.err = msg.err
		m.done = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m spinnerModel) View() string {
	if m.done {
		return ""
	}
	return "  " + m.spinner.View() + " " + mutedStyle.Render(m.message) + "\n"
}

// Spinner displays an animated spinner while fn executes.
// Returns the error from fn, if any.
//
// Without an interactive terminal (CI, headless, piped output) the
// bubbletea program cannot open /dev/tty and p.Run would fail, masking fn's
// result entirely. In that case Spinner degrades to a plain one-line message
// and runs fn directly, so the wrapped work still executes and its real error
// is returned unchanged.
func Spinner(message string, fn func() error) error {
	if !interactive() {
		fprintf("  %s\n", message)
		return fn()
	}

	m := newSpinnerModel(message)
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

	go func() {
		err := fn()
		p.Send(spinnerDoneMsg{err: err})
	}()

	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	if fm, ok := finalModel.(spinnerModel); ok {
		return fm.err
	}
	return nil
}

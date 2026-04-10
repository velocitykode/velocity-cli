package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type searchModel struct {
	label     string
	textInput textinput.Model
	searchFn  func(string) []string
	results   []string
	cursor    int
	done      bool
}

func newSearchModel(label string, fn func(string) []string) searchModel {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40
	ti.PromptStyle = cursorStyle
	ti.TextStyle = boldStyle

	// Show initial results with empty query
	results := fn("")

	return searchModel{
		label:     label,
		textInput: ti,
		searchFn:  fn,
		results:   results,
	}
}

func (m searchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case tea.KeyDown:
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	prev := m.textInput.Value()
	m.textInput, cmd = m.textInput.Update(msg)

	// Re-search on text change
	if m.textInput.Value() != prev {
		m.results = m.searchFn(m.textInput.Value())
		m.cursor = 0
	}

	return m, cmd
}

func (m searchModel) View() string {
	if m.done {
		val := m.value()
		if val == "" {
			return fmt.Sprintf("  %s %s\n", promptLabelStyle.Render(m.label), mutedStyle.Render("(none)")) + "\n"
		}
		return fmt.Sprintf("  %s %s\n", promptLabelStyle.Render(m.label), successStyle.Render(val)) + "\n"
	}

	s := fmt.Sprintf("  %s\n  %s\n\n", promptLabelStyle.Render(m.label), m.textInput.View())

	maxShow := 7
	for i, r := range m.results {
		if i >= maxShow {
			s += fmt.Sprintf("  %s\n", mutedStyle.Render(fmt.Sprintf("  ... and %d more", len(m.results)-maxShow)))
			break
		}
		if i == m.cursor {
			s += fmt.Sprintf("  %s %s\n", cursorStyle.Render("▸"), selectedStyle.Render(r))
		} else {
			s += fmt.Sprintf("    %s\n", unselectedStyle.Render(r))
		}
	}

	if len(m.results) == 0 {
		s += fmt.Sprintf("  %s\n", mutedStyle.Render("  No results"))
	}

	s += fmt.Sprintf("\n  %s\n", promptHelpStyle.Render("↑/↓ navigate, enter select"))
	return s
}

func (m searchModel) value() string {
	if m.cursor >= 0 && m.cursor < len(m.results) {
		return m.results[m.cursor]
	}
	return ""
}

// Search displays a text input that filters results dynamically.
// The fn is called with the current query and should return matching options.
func Search(label string, fn func(query string) []string) string {
	m := newSearchModel(label, fn)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return ""
	}
	if fm, ok := finalModel.(searchModel); ok {
		return fm.value()
	}
	return ""
}

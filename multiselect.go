package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type multiselectModel struct {
	label    string
	options  []string
	cursor   int
	selected map[int]bool
	done     bool
}

func newMultiselectModel(label string, options []string) multiselectModel {
	return multiselectModel{
		label:    label,
		options:  options,
		selected: make(map[int]bool),
	}
}

func (m multiselectModel) Init() tea.Cmd { return nil }

func (m multiselectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.selected = nil
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
		case " ":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "a":
			allSelected := true
			for i := range m.options {
				if !m.selected[i] {
					allSelected = false
					break
				}
			}
			for i := range m.options {
				m.selected[i] = !allSelected
			}
		}
	}
	return m, nil
}

func (m multiselectModel) View() string {
	if m.done {
		vals := m.values()
		if len(vals) == 0 {
			return fmt.Sprintf("  %s %s\n", promptLabelStyle.Render(m.label), mutedStyle.Render("(none)")) + "\n"
		}
		s := fmt.Sprintf("  %s\n", promptLabelStyle.Render(m.label))
		for _, v := range vals {
			s += fmt.Sprintf("  %s %s\n", symbolCheck, successStyle.Render(v))
		}
		return s + "\n"
	}

	s := fmt.Sprintf("  %s\n\n", promptLabelStyle.Render(m.label))
	for i, opt := range m.options {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("▸ ")
		}

		check := unselectedStyle.Render("○")
		label := unselectedStyle.Render(opt)
		if m.selected[i] {
			check = selectedStyle.Render("●")
			label = selectedStyle.Render(opt)
		}

		s += fmt.Sprintf("  %s%s %s\n", cursor, check, label)
	}
	s += fmt.Sprintf("\n  %s\n", promptHelpStyle.Render("↑/↓ navigate, space toggle, a toggle all, enter confirm"))
	return s
}

func (m multiselectModel) values() []string {
	var result []string
	for i, opt := range m.options {
		if m.selected[i] {
			result = append(result, opt)
		}
	}
	return result
}

// Multiselect displays a list of options with toggles and returns all selected options.
func Multiselect(label string, options []string) []string {
	m := newMultiselectModel(label, options)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil
	}
	if fm, ok := finalModel.(multiselectModel); ok {
		return fm.values()
	}
	return nil
}

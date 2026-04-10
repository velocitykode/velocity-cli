package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type progressIncrMsg struct{}

type progressModel struct {
	progress  progress.Model
	total     int
	completed int
}

func newProgressModel(total int) progressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)
	return progressModel{progress: p, total: total}
}

func (m progressModel) Init() tea.Cmd {
	return nil
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case progressIncrMsg:
		m.completed++
		percent := float64(m.completed) / float64(m.total)
		if percent >= 1.0 {
			return m, tea.Sequence(
				m.progress.SetPercent(1.0),
				tea.Quit,
			)
		}
		return m, m.progress.SetPercent(percent)
	case progress.FrameMsg:
		model, cmd := m.progress.Update(msg)
		m.progress = model.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m progressModel) View() string {
	return fmt.Sprintf("\n  %s  %d/%d\n\n", m.progress.View(), m.completed, m.total)
}

// Progress displays an animated progress bar for batch operations.
// The fn receives an increment function to call after each completed item.
func Progress(total int, fn func(increment func())) {
	if total <= 0 {
		return
	}

	m := newProgressModel(total)
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

	go func() {
		fn(func() {
			p.Send(progressIncrMsg{})
		})
	}()

	p.Run()
}

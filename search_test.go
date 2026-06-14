package prism

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func filterFn(query string) []string {
	all := []string{"Alice", "Bob", "Charlie", "Dave"}
	if query == "" {
		return all
	}
	var results []string
	for _, s := range all {
		if strings.Contains(strings.ToLower(s), strings.ToLower(query)) {
			results = append(results, s)
		}
	}
	return results
}

func TestSearchModel_InitialResults(t *testing.T) {
	m := newSearchModel("Find:", filterFn)
	if len(m.results) != 4 {
		t.Errorf("expected 4 initial results, got %d", len(m.results))
	}
}

func TestSearchModel_Navigation(t *testing.T) {
	m := newSearchModel("Find:", filterFn)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	sm := updated.(searchModel)
	if sm.cursor != 1 {
		t.Errorf("expected cursor=1, got %d", sm.cursor)
	}

	updated, _ = sm.Update(tea.KeyMsg{Type: tea.KeyUp})
	sm = updated.(searchModel)
	if sm.cursor != 0 {
		t.Errorf("expected cursor=0, got %d", sm.cursor)
	}
}

func TestSearchModel_Enter(t *testing.T) {
	m := newSearchModel("Find:", filterFn)
	m.cursor = 2
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	sm := updated.(searchModel)
	if !sm.done {
		t.Error("enter should set done")
	}
	if sm.value() != "Charlie" {
		t.Errorf("expected 'Charlie', got %q", sm.value())
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestSearchModel_CtrlC(t *testing.T) {
	m := newSearchModel("Find:", filterFn)
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	sm := updated.(searchModel)
	if !sm.Cancelled() {
		t.Error("ctrl+c should mark model cancelled")
	}
	if !sm.done {
		t.Error("ctrl+c should set done")
	}
	if sm.value() != "" {
		t.Error("ctrl+c should return empty value")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestSearchModel_View(t *testing.T) {
	m := newSearchModel("Find user:", filterFn)
	v := m.View()
	if !strings.Contains(v, "Find user:") {
		t.Error("view should show label")
	}
	if !strings.Contains(v, "Alice") {
		t.Error("view should show results")
	}
}

func TestSearchModel_MaxResults(t *testing.T) {
	many := func(q string) []string {
		results := make([]string, 20)
		for i := range results {
			results[i] = strings.Repeat("x", i+1)
		}
		return results
	}
	m := newSearchModel("Find:", many)
	v := m.View()
	if !strings.Contains(v, "and 13 more") {
		t.Error("should show overflow count for >7 results")
	}
}

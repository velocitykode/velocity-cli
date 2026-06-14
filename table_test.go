package prism

import (
	"strings"
	"testing"
)

func TestTable_Basic(t *testing.T) {
	out := capture(func() {
		Table(
			[]string{"Method", "Path", "Name"},
			[][]string{
				{"GET", "/users", "users.index"},
				{"POST", "/users", "users.store"},
			},
		)
	})

	for _, want := range []string{"Method", "Path", "Name", "GET", "/users", "users.index", "POST", "users.store", "Showing 2 rows"} {
		if !strings.Contains(out, want) {
			t.Errorf("Table() missing %q in output", want)
		}
	}
}

func TestTable_EmptyHeaders(t *testing.T) {
	out := capture(func() {
		Table([]string{}, [][]string{})
	})
	if out != "" {
		t.Errorf("Table() with no headers should produce no output, got: %q", out)
	}
}

func TestTable_EmptyRows(t *testing.T) {
	out := capture(func() {
		Table([]string{"A", "B"}, [][]string{})
	})
	if !strings.Contains(out, "A") || !strings.Contains(out, "B") {
		t.Error("Table() should still render headers with no rows")
	}
	if !strings.Contains(out, "Showing 0 rows") {
		t.Error("Table() should show 0 rows count")
	}
}

func TestTable_UnevenRows(t *testing.T) {
	out := capture(func() {
		Table(
			[]string{"A", "B", "C"},
			[][]string{
				{"1"},
				{"1", "2", "3"},
			},
		)
	})
	if !strings.Contains(out, "1") && !strings.Contains(out, "3") {
		t.Error("Table() should handle uneven rows")
	}
}

package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

// Table prints a formatted table with styled headers and dynamic column widths.
func Table(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	cols := len(headers)

	// Calculate max width per column
	widths := make([]int, cols)
	for i, h := range headers {
		widths[i] = lipgloss.Width(h)
	}
	for _, row := range rows {
		for i := 0; i < cols && i < len(row); i++ {
			w := lipgloss.Width(row[i])
			if w > widths[i] {
				widths[i] = w
			}
		}
	}

	// Add padding
	for i := range widths {
		widths[i] += 2
	}

	// Render header
	Newline()
	var headerParts []string
	for i, h := range headers {
		headerParts = append(headerParts, tableHeaderStyle.Render(padRight(h, widths[i])))
	}
	fprintf("  %s\n", strings.Join(headerParts, " "))

	// Separator
	totalWidth := 0
	for _, w := range widths {
		totalWidth += w + 1
	}
	fprintf("  %s\n", tableSepStyle.Render(strings.Repeat("─", totalWidth)))

	// Render rows
	for _, row := range rows {
		var parts []string
		for i := 0; i < cols; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			parts = append(parts, padRight(cell, widths[i]))
		}
		fprintf("  %s\n", strings.Join(parts, " "))
	}

	Newline()
	Muted(fmt.Sprintf("Showing %d rows", len(rows)))
	Newline()
}

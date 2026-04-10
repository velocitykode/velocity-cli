package cli

import (
	"fmt"
	"strings"
)

// Header prints a styled command header.
func Header(message string) {
	fprintf("\n%s\n\n", primaryStyle.Render(strings.ToUpper(message)))
}

// Info prints an informational message with an arrow symbol.
func Info(message string) {
	fprintf("%s %s\n", symbolArrow, mutedStyle.Render(message))
}

// Success prints a success message with a checkmark.
func Success(message string) {
	fprintf("%s %s\n", symbolCheck, successStyle.Render(message))
}

// Warning prints a warning message.
func Warning(message string) {
	fprintf("%s %s\n", symbolWarn, warningStyle.Render(message))
}

// Error prints an error message.
func Error(message string) {
	fprintf("%s %s\n", symbolCross, errorStyle.Render(message))
}

// Muted prints dimmed text.
func Muted(message string) {
	fprintf("  %s\n", mutedStyle.Render(message))
}

// Bold prints bold text.
func Bold(message string) {
	fprintln(boldStyle.Render(message))
}

// Newline prints an empty line.
func Newline() {
	fprintln()
}

// Highlight returns styled text without printing.
func Highlight(text string) string {
	return primaryStyle.Render(text)
}

// KeyValue prints a key-value pair.
func KeyValue(key, value string) {
	fprintf("  %s %s\n", mutedStyle.Render(key+":"), value)
}

// Step prints an indented step message.
func Step(message string) {
	fprintf("  %s\n", mutedStyle.Render(message))
}

// NextSteps prints formatted next steps.
func NextSteps(steps []string) {
	Newline()
	fprintln(mutedStyle.Render("Next steps:"))
	for i, step := range steps {
		fprintf("  %s %s\n", primaryStyle.Render(fmt.Sprintf("%d.", i+1)), step)
	}
}

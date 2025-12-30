package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor = lipgloss.Color("#0e87cd")
	successColor = lipgloss.Color("#10b981")
	warningColor = lipgloss.Color("#f59e0b")
	errorColor   = lipgloss.Color("#ef4444")
	mutedColor   = lipgloss.Color("#6b7280")

	// Symbols
	arrowSymbol = lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("→")
	checkSymbol = lipgloss.NewStyle().Foreground(successColor).Render("✓")
	warnSymbol  = lipgloss.NewStyle().Foreground(warningColor).Render("!")
	crossSymbol = lipgloss.NewStyle().Foreground(errorColor).Render("✗")

	// Text styles
	mutedStyle   = lipgloss.NewStyle().Foreground(mutedColor)
	primaryStyle = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(warningColor).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
)

// Header prints a styled command header (uppercase cyan)
func Header(command string) {
	fmt.Printf("\n%s\n\n", primaryStyle.Render(strings.ToUpper(command)))
}

// Info prints an info message with arrow symbol
func Info(message string) {
	fmt.Printf("%s %s\n", arrowSymbol, successStyle.Render(message))
}

// Success prints a success message with checkmark
func Success(message string) {
	fmt.Printf("%s %s\n", checkSymbol, successStyle.Render(message))
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Printf("%s %s\n", warnSymbol, warningStyle.Render(message))
}

// Error prints an error message
func Error(message string) {
	fmt.Printf("%s %s\n", crossSymbol, errorStyle.Render(message))
}

// Step prints a muted step message (indented)
func Step(message string) {
	fmt.Printf("  %s\n", mutedStyle.Render(message))
}

// Muted prints muted text
func Muted(message string) {
	fmt.Printf("  %s\n", mutedStyle.Render(message))
}

// Bold prints bold text
func Bold(message string) {
	fmt.Println(lipgloss.NewStyle().Bold(true).Render(message))
}

// Highlight returns highlighted text
func Highlight(text string) string {
	return primaryStyle.Render(text)
}

// Command returns styled command text
func Command(cmd string) string {
	return primaryStyle.Render(cmd)
}

// KeyValue prints a key-value pair
func KeyValue(key, value string) {
	fmt.Printf("  %s %s\n", mutedStyle.Render(key+":"), value)
}

// Newline prints an empty line
func Newline() {
	fmt.Println()
}

// NextSteps prints formatted next steps
func NextSteps(steps []string) {
	fmt.Println()
	fmt.Println(mutedStyle.Render("Next steps:"))
	for i, step := range steps {
		fmt.Printf("  %s %s\n", primaryStyle.Render(fmt.Sprintf("%d.", i+1)), step)
	}
}

// Task runs an action with step message, then shows success/error
func Task(stepMsg, successMsg string, action func() error) error {
	Info(stepMsg)
	fmt.Printf("  %s\n", mutedStyle.Render(stepMsg+"..."))
	if err := action(); err != nil {
		return err
	}
	fmt.Println()
	Success(successMsg)
	return nil
}

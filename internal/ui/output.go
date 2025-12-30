package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Brand colors
	primaryColor = lipgloss.Color("#0e87cd")
	successColor = lipgloss.Color("#10b981")
	warningColor = lipgloss.Color("#f59e0b")
	errorColor   = lipgloss.Color("#ef4444")
	mutedColor   = lipgloss.Color("#6b7280")

	// Styles
	primaryStyle = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(warningColor).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
	mutedStyle   = lipgloss.NewStyle().Foreground(mutedColor)
	boldStyle    = lipgloss.NewStyle().Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1)
)

// Header prints a styled command header
func Header(command string) {
	fmt.Println()
	fmt.Println(primaryStyle.Render("VELOCITY " + strings.ToUpper(command)))
	fmt.Println()
}

// Step prints a progress step
func Step(message string) {
	fmt.Println(mutedStyle.Render("  " + message))
}

// Success prints a success message
func Success(message string) {
	fmt.Println(successStyle.Render("✓ " + message))
}

// SuccessBox prints a success message in a box
func SuccessBox(message string) {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(successColor).
		Padding(0, 1)
	fmt.Println(style.Render(successStyle.Render("✓ ") + message))
}

// Error prints an error message
func Error(message string) {
	fmt.Println(errorStyle.Render("✗ " + message))
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Println(warningStyle.Render("⚠ " + message))
}

// Info prints an info message
func Info(message string) {
	fmt.Println(primaryStyle.Render("→ " + message))
}

// Muted prints a muted/secondary message
func Muted(message string) {
	fmt.Println(mutedStyle.Render("  " + message))
}

// Bold prints bold text
func Bold(message string) {
	fmt.Println(boldStyle.Render(message))
}

// NextSteps prints a formatted next steps box
func NextSteps(steps []string) {
	fmt.Println()
	fmt.Println(mutedStyle.Render("Next steps:"))
	for i, step := range steps {
		fmt.Printf("  %s %s\n", primaryStyle.Render(fmt.Sprintf("%d.", i+1)), step)
	}
	fmt.Println()
}

// Command prints a command suggestion
func Command(cmd string) string {
	return primaryStyle.Render(cmd)
}

// Highlight returns highlighted text
func Highlight(text string) string {
	return primaryStyle.Render(text)
}

// Box prints content in a bordered box
func Box(content string) {
	fmt.Println(boxStyle.Render(content))
}

// KeyValue prints a key-value pair
func KeyValue(key, value string) {
	fmt.Printf("  %s %s\n", mutedStyle.Render(key+":"), value)
}

// Divider prints a subtle divider line
func Divider() {
	fmt.Println(mutedStyle.Render(strings.Repeat("─", 40)))
}

// Newline prints an empty line
func Newline() {
	fmt.Println()
}

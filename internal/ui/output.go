package ui

import (
	"fmt"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor = lipgloss.Color("#0e87cd")
	successColor = lipgloss.Color("#10b981")
	warningColor = lipgloss.Color("#f59e0b")
	errorColor   = lipgloss.Color("#ef4444")
	mutedColor   = lipgloss.Color("#6b7280")

	// Label styles
	labelStyle   = lipgloss.NewStyle().Bold(true).Padding(0, 1)
	infoLabel    = labelStyle.Background(primaryColor).Foreground(lipgloss.Color("#ffffff")).Render("INFO")
	successLabel = labelStyle.Background(successColor).Foreground(lipgloss.Color("#ffffff")).Render("SUCCESS")
	warningLabel = labelStyle.Background(warningColor).Foreground(lipgloss.Color("#000000")).Render("WARNING")
	errorLabel   = labelStyle.Background(errorColor).Foreground(lipgloss.Color("#ffffff")).Render("ERROR")

	// Text styles
	mutedStyle   = lipgloss.NewStyle().Foreground(mutedColor)
	primaryStyle = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
)

// Header prints a styled command header
func Header(command string) {
	fmt.Printf("%s %s\n", infoLabel, command)
}

// Info prints an info message with label
func Info(message string) {
	fmt.Printf("%s %s\n", infoLabel, message)
}

// Success prints a success message with label
func Success(message string) {
	fmt.Printf("%s %s\n", successLabel, message)
}

// Warning prints a warning message with label
func Warning(message string) {
	fmt.Printf("%s %s\n", warningLabel, message)
}

// Error prints an error message with label
func Error(message string) {
	fmt.Printf("%s %s\n", errorLabel, message)
}

// Step prints a muted step message (no label)
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

// Spinner runs an action with a spinner
func Spinner(title string, action func()) error {
	return spinner.New().
		Title(title).
		Action(action).
		Run()
}

// SpinnerWithError runs an action with a spinner and returns error
func SpinnerWithError(title string, action func() error) error {
	var err error
	spinErr := spinner.New().
		Title(title).
		Action(func() {
			err = action()
		}).
		Run()
	if spinErr != nil {
		return spinErr
	}
	return err
}

package cli

import "github.com/charmbracelet/lipgloss"

// Brand colors
var (
	colorPrimary = lipgloss.Color("#0e87cd")
	colorSuccess = lipgloss.Color("#10b981")
	colorWarning = lipgloss.Color("#f59e0b")
	colorError   = lipgloss.Color("#ef4444")
	colorMuted   = lipgloss.Color("#6b7280")
)

// Text styles
var (
	primaryStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(colorError).Bold(true)
	mutedStyle   = lipgloss.NewStyle().Foreground(colorMuted)
	boldStyle    = lipgloss.NewStyle().Bold(true)
)

// Symbols
var (
	symbolArrow = primaryStyle.Render("→")
	symbolCheck = lipgloss.NewStyle().Foreground(colorSuccess).Render("✓")
	symbolWarn  = lipgloss.NewStyle().Foreground(colorWarning).Render("!")
	symbolCross = lipgloss.NewStyle().Foreground(colorError).Render("✗")
)

// Box styles for Note/Alert
var (
	noteBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 2)

	alertBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorError).
			Padding(0, 2)
)

// Table styles
var (
	tableHeaderStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	tableSepStyle    = lipgloss.NewStyle().Foreground(colorMuted)
)

// Prompt styles
var (
	promptLabelStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	promptHelpStyle  = lipgloss.NewStyle().Foreground(colorMuted)
	cursorStyle      = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	selectedStyle    = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	unselectedStyle  = lipgloss.NewStyle().Foreground(colorMuted)
	errorMsgStyle    = lipgloss.NewStyle().Foreground(colorError)
)

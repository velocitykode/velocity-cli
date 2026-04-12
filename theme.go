package cli

import "github.com/charmbracelet/lipgloss"

// Colors in use, rebuilt whenever Configure is called.
var (
	colorPrimary lipgloss.Color
	colorSuccess lipgloss.Color
	colorWarning lipgloss.Color
	colorError   lipgloss.Color
	colorMuted   lipgloss.Color
)

// Text styles.
var (
	primaryStyle lipgloss.Style
	successStyle lipgloss.Style
	warningStyle lipgloss.Style
	errorStyle   lipgloss.Style
	mutedStyle   lipgloss.Style
	boldStyle    lipgloss.Style
)

// Pre-rendered symbols.
var (
	symbolArrow string
	symbolCheck string
	symbolWarn  string
	symbolCross string
)

// Box styles for Note/Alert.
var (
	noteBoxStyle  lipgloss.Style
	alertBoxStyle lipgloss.Style
)

// Table styles.
var (
	tableHeaderStyle lipgloss.Style
	tableSepStyle    lipgloss.Style
)

// Prompt styles.
var (
	promptLabelStyle lipgloss.Style
	promptHelpStyle  lipgloss.Style
	cursorStyle      lipgloss.Style
	selectedStyle    lipgloss.Style
	unselectedStyle  lipgloss.Style
	errorMsgStyle    lipgloss.Style
)

func init() {
	applyTheme(defaultConfig())
}

func applyTheme(cfg Config) {
	colorPrimary = lipgloss.Color(cfg.Colors.Primary)
	colorSuccess = lipgloss.Color(cfg.Colors.Success)
	colorWarning = lipgloss.Color(cfg.Colors.Warning)
	colorError = lipgloss.Color(cfg.Colors.Error)
	colorMuted = lipgloss.Color(cfg.Colors.Muted)

	primaryStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
	errorStyle = lipgloss.NewStyle().Foreground(colorError).Bold(true)
	mutedStyle = lipgloss.NewStyle().Foreground(colorMuted)
	boldStyle = lipgloss.NewStyle().Bold(true)

	symbolArrow = primaryStyle.Render(cfg.Symbols.Arrow)
	symbolCheck = lipgloss.NewStyle().Foreground(colorSuccess).Render(cfg.Symbols.Check)
	symbolWarn = lipgloss.NewStyle().Foreground(colorWarning).Render(cfg.Symbols.Warn)
	symbolCross = lipgloss.NewStyle().Foreground(colorError).Render(cfg.Symbols.Cross)

	noteBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(0, 2)

	alertBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorError).
		Padding(0, 2)

	tableHeaderStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	tableSepStyle = lipgloss.NewStyle().Foreground(colorMuted)

	promptLabelStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	promptHelpStyle = lipgloss.NewStyle().Foreground(colorMuted)
	cursorStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true)
	unselectedStyle = lipgloss.NewStyle().Foreground(colorMuted)
	errorMsgStyle = lipgloss.NewStyle().Foreground(colorError)

	activeConfig = cfg
}

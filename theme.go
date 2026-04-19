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
	// "default" is a sentinel meaning "do not set a foreground color" -
	// primary-styled text renders bold in the terminal's own default fg,
	// useful for CLIs that want to reserve green for success, red for
	// error, etc. Any other value (including "") is passed to lipgloss
	// as-is.
	primaryIsDefault := cfg.Colors.Primary == "default"
	if primaryIsDefault {
		colorPrimary = lipgloss.Color("")
	} else {
		colorPrimary = lipgloss.Color(cfg.Colors.Primary)
	}
	colorSuccess = lipgloss.Color(cfg.Colors.Success)
	colorWarning = lipgloss.Color(cfg.Colors.Warning)
	colorError = lipgloss.Color(cfg.Colors.Error)
	colorMuted = lipgloss.Color(cfg.Colors.Muted)

	primaryStyle = lipgloss.NewStyle().Bold(true)
	if !primaryIsDefault {
		primaryStyle = primaryStyle.Foreground(colorPrimary)
	}
	successStyle = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
	errorStyle = lipgloss.NewStyle().Foreground(colorError).Bold(true)
	mutedStyle = lipgloss.NewStyle().Foreground(colorMuted)
	boldStyle = lipgloss.NewStyle().Bold(true)

	symbolArrow = primaryStyle.Render(cfg.Symbols.Arrow)
	symbolCheck = lipgloss.NewStyle().Foreground(colorSuccess).Render(cfg.Symbols.Check)
	symbolWarn = lipgloss.NewStyle().Foreground(colorWarning).Render(cfg.Symbols.Warn)
	symbolCross = lipgloss.NewStyle().Foreground(colorError).Render(cfg.Symbols.Cross)

	// Helper: apply Foreground(colorPrimary) only when primary isn't the
	// "default" sentinel. Keeps lipgloss from emitting an empty-color ANSI.
	withPrimaryFg := func(s lipgloss.Style) lipgloss.Style {
		if primaryIsDefault {
			return s
		}
		return s.Foreground(colorPrimary)
	}

	noteBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 2)
	if !primaryIsDefault {
		noteBox = noteBox.BorderForeground(colorPrimary)
	}
	noteBoxStyle = noteBox

	alertBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorError).
		Padding(0, 2)

	tableHeaderStyle = withPrimaryFg(lipgloss.NewStyle().Bold(true))
	tableSepStyle = lipgloss.NewStyle().Foreground(colorMuted)

	promptLabelStyle = withPrimaryFg(lipgloss.NewStyle().Bold(true))
	promptHelpStyle = lipgloss.NewStyle().Foreground(colorMuted)
	cursorStyle = withPrimaryFg(lipgloss.NewStyle().Bold(true))
	selectedStyle = withPrimaryFg(lipgloss.NewStyle().Bold(true))
	unselectedStyle = lipgloss.NewStyle().Foreground(colorMuted)
	errorMsgStyle = lipgloss.NewStyle().Foreground(colorError)

	activeConfig = cfg
}

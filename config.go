package cli

import (
	"io"

	"github.com/BurntSushi/toml"
)

// Config controls colors and symbols used by all output and prompt components.
// Empty string fields are ignored by Configure - they keep the current value.
type Config struct {
	Colors  Colors  `toml:"colors"`
	Symbols Symbols `toml:"symbols"`
}

// Colors are hex strings like "#0e87cd" or named ANSI colors accepted by lipgloss.
type Colors struct {
	Primary string `toml:"primary"`
	Success string `toml:"success"`
	Warning string `toml:"warning"`
	Error   string `toml:"error"`
	Muted   string `toml:"muted"`
}

// Symbols are the glyphs used by Info (arrow), Success (check),
// Warning (warn), Error (cross), and Tip (info glyph).
type Symbols struct {
	Arrow string `toml:"arrow"`
	Check string `toml:"check"`
	Warn  string `toml:"warn"`
	Cross string `toml:"cross"`
	Tip   string `toml:"tip"`
}

var activeConfig Config

func defaultConfig() Config {
	return Config{
		Colors: Colors{
			Primary: "#0e87cd",
			Success: "#10b981",
			Warning: "#f59e0b",
			Error:   "#ef4444",
			Muted:   "#6b7280",
		},
		Symbols: Symbols{
			Arrow: "→",
			Check: "✓",
			Warn:  "!",
			Cross: "✗",
			Tip:   "ℹ",
		},
	}
}

// ActiveConfig returns the currently applied config.
func ActiveConfig() Config {
	return activeConfig
}

// Configure merges cfg into the active config and rebuilds all styles.
// Empty string fields are ignored (they keep the current value), so callers
// can override just the fields they care about.
func Configure(cfg Config) {
	merged := activeConfig
	if cfg.Colors.Primary != "" {
		merged.Colors.Primary = cfg.Colors.Primary
	}
	if cfg.Colors.Success != "" {
		merged.Colors.Success = cfg.Colors.Success
	}
	if cfg.Colors.Warning != "" {
		merged.Colors.Warning = cfg.Colors.Warning
	}
	if cfg.Colors.Error != "" {
		merged.Colors.Error = cfg.Colors.Error
	}
	if cfg.Colors.Muted != "" {
		merged.Colors.Muted = cfg.Colors.Muted
	}
	if cfg.Symbols.Arrow != "" {
		merged.Symbols.Arrow = cfg.Symbols.Arrow
	}
	if cfg.Symbols.Check != "" {
		merged.Symbols.Check = cfg.Symbols.Check
	}
	if cfg.Symbols.Warn != "" {
		merged.Symbols.Warn = cfg.Symbols.Warn
	}
	if cfg.Symbols.Cross != "" {
		merged.Symbols.Cross = cfg.Symbols.Cross
	}
	if cfg.Symbols.Tip != "" {
		merged.Symbols.Tip = cfg.Symbols.Tip
	}
	applyTheme(merged)
}

// LoadConfig reads TOML config from r and applies it. Missing fields keep defaults.
// The caller owns the file - open with os.Open, pass the *os.File, close it yourself.
func LoadConfig(r io.Reader) error {
	var cfg Config
	if _, err := toml.NewDecoder(r).Decode(&cfg); err != nil {
		return err
	}
	Configure(cfg)
	return nil
}

// ResetConfig reverts all styles to the built-in defaults.
func ResetConfig() {
	applyTheme(defaultConfig())
}

// StylePrimary renders text in the primary brand color (same as Highlight).
func StylePrimary(text string) string { return primaryStyle.Render(text) }

// StyleMuted renders text in the muted color.
func StyleMuted(text string) string { return mutedStyle.Render(text) }

// StyleSuccess renders text in the success color.
func StyleSuccess(text string) string { return successStyle.Render(text) }

// StyleWarning renders text in the warning color.
func StyleWarning(text string) string { return warningStyle.Render(text) }

// StyleError renders text in the error color.
func StyleError(text string) string { return errorStyle.Render(text) }

package cli

import (
	"strings"
	"testing"
)

func TestDefaultConfig_Applied(t *testing.T) {
	ResetConfig()
	cfg := ActiveConfig()
	if cfg.Colors.Primary != "#0e87cd" {
		t.Errorf("expected default primary color, got %q", cfg.Colors.Primary)
	}
	if cfg.Symbols.Check != "✓" {
		t.Errorf("expected default check symbol, got %q", cfg.Symbols.Check)
	}
}

func TestConfigure_OverridesColors(t *testing.T) {
	defer ResetConfig()
	Configure(Config{Colors: Colors{Primary: "#ff00ff"}})
	if got := ActiveConfig().Colors.Primary; got != "#ff00ff" {
		t.Errorf("expected #ff00ff, got %q", got)
	}
	if got := ActiveConfig().Colors.Muted; got != "#6b7280" {
		t.Errorf("unset field should keep default, got %q", got)
	}
}

func TestConfigure_OverridesSymbols(t *testing.T) {
	defer ResetConfig()
	Configure(Config{Symbols: Symbols{Check: "OK"}})
	if got := ActiveConfig().Symbols.Check; got != "OK" {
		t.Errorf("expected OK, got %q", got)
	}
	if got := ActiveConfig().Symbols.Arrow; got != "→" {
		t.Errorf("unset symbol should keep default, got %q", got)
	}
}

func TestLoadConfig_FromReader(t *testing.T) {
	defer ResetConfig()
	body := `
[colors]
primary = "#abcdef"

[symbols]
arrow = ">>"
`
	if err := LoadConfig(strings.NewReader(body)); err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if got := ActiveConfig().Colors.Primary; got != "#abcdef" {
		t.Errorf("expected #abcdef, got %q", got)
	}
	if got := ActiveConfig().Symbols.Arrow; got != ">>" {
		t.Errorf("expected >>, got %q", got)
	}
}

func TestLoadConfig_InvalidTOML(t *testing.T) {
	defer ResetConfig()
	if err := LoadConfig(strings.NewReader("not = [valid")); err == nil {
		t.Fatal("expected error for invalid TOML")
	}
}

func TestStyleHelpers_RenderSomething(t *testing.T) {
	ResetConfig()
	for name, fn := range map[string]func(string) string{
		"primary": StylePrimary,
		"muted":   StyleMuted,
		"success": StyleSuccess,
		"warning": StyleWarning,
		"error":   StyleError,
	} {
		out := fn("hello")
		if !strings.Contains(out, "hello") {
			t.Errorf("%s helper dropped text: %q", name, out)
		}
	}
}

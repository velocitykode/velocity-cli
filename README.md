# Prism

**Build beautiful command-line tools in Go, without the boilerplate.** A component library for consistent output styling, tables, spinners, progress bars, and interactive prompts. Import it into any Go project and ship a polished CLI in an afternoon.

```bash
go get github.com/velocitykode/prism
```

```go
import "github.com/velocitykode/prism"
```

> Prism has **zero framework dependency**. Use it in any Go program, Velocity or not. It is built on the [Charm](https://charm.sh) stack (lipgloss, bubbletea, bubbles).

## Why Prism

- **Consistent by default.** Leveled output, boxes, and prompts share one theme, so your whole tool looks intentional.
- **Batteries included.** Output, tables, spinners, progress bars, and seven prompt types in one import.
- **One-liners.** Each helper is a single call that returns a typed result; no models, no event loops to wire up.
- **Drop-in anywhere.** No framework lock-in, no global state.

## Output

```go
prism.Header("migrate")               // MIGRATE (styled header)
prism.Info("Running migrations...")   // informational line
prism.Success("Done")                 // checkmark, success color
prism.Warning("No migrations found")  // warning glyph and color
prism.Error("Connection failed")      // error glyph and color
prism.Muted("skipping...")            // dimmed text
prism.Bold("important")               // bold text

prism.Note("Server running on port 4000")  // boxed message (primary border)
prism.Alert("Database connection failed")  // boxed message (error border)
```

## Tables

```go
prism.Table(
    []string{"Method", "Path", "Name"},
    [][]string{
        {"GET", "/users", "users.index"},
        {"POST", "/users", "users.store"},
    },
)
```

## Spinners & Progress

A spinner wraps a slow step; a progress bar tracks a known count:

```go
prism.Spinner("Building...", func() error {
    return exec.Command("go", "build", ".").Run()
})

prism.Progress(len(items), func(inc func()) {
    for _, item := range items {
        process(item)
        inc()
    }
})
```

## Interactive Prompts

Each prompt blocks for input and returns the typed result:

```go
name := prism.Text("Model name:", prism.WithRequired(), prism.WithPlaceholder("User"))
pass := prism.Password("Database password:")
yes  := prism.Confirm("Generate migration?", prism.WithDefaultYes())
db   := prism.Select("Database driver:", []string{"mysql", "postgres", "sqlite"})
features := prism.Multiselect("Enable:", []string{"soft-deletes", "uuid", "timestamps"})

user := prism.Search("Find user:", func(q string) []string {
    return filterUsers(q)
})
```

## Theming

Output colors and prompt styling are driven by a shared theme, so a whole tool stays visually consistent. Override it to match your brand.

## Built On

- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - styling
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - interactive prompts
- [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) - prebuilt components

## License

MIT

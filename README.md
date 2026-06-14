# Velocity CLI

**Build beautiful command-line tools in Go, without the boilerplate.** A component library for consistent output styling, tables, spinners, progress bars, and interactive prompts. Import it into any Go project and ship a polished CLI in an afternoon.

```bash
go get github.com/velocitykode/velocity-cli
```

```go
import cli "github.com/velocitykode/velocity-cli"
```

> Despite the name, Velocity CLI has **zero framework dependency**. Use it in any Go program, Velocity or not. It is built on the [Charm](https://charm.sh) stack (lipgloss, bubbletea, bubbles).

## Why Velocity CLI

- **Consistent by default.** Leveled output, boxes, and prompts share one theme, so your whole tool looks intentional.
- **Batteries included.** Output, tables, spinners, progress bars, and seven prompt types in one import.
- **One-liners.** Each helper is a single call that returns a typed result; no models, no event loops to wire up.
- **Drop-in anywhere.** No framework lock-in, no global state.

## Output

```go
cli.Header("migrate")               // MIGRATE (styled header)
cli.Info("Running migrations...")   // informational line
cli.Success("Done")                 // checkmark, success color
cli.Warning("No migrations found")  // warning glyph and color
cli.Error("Connection failed")      // error glyph and color
cli.Muted("skipping...")            // dimmed text
cli.Bold("important")               // bold text

cli.Note("Server running on port 4000")  // boxed message (primary border)
cli.Alert("Database connection failed")  // boxed message (error border)
```

## Tables

```go
cli.Table(
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
cli.Spinner("Building...", func() error {
    return exec.Command("go", "build", ".").Run()
})

cli.Progress(len(items), func(inc func()) {
    for _, item := range items {
        process(item)
        inc()
    }
})
```

## Interactive Prompts

Each prompt blocks for input and returns the typed result:

```go
name := cli.Text("Model name:", cli.WithRequired(), cli.WithPlaceholder("User"))
pass := cli.Password("Database password:")
yes  := cli.Confirm("Generate migration?", cli.WithDefaultYes())
db   := cli.Select("Database driver:", []string{"mysql", "postgres", "sqlite"})
features := cli.Multiselect("Enable:", []string{"soft-deletes", "uuid", "timestamps"})

user := cli.Search("Find user:", func(q string) []string {
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

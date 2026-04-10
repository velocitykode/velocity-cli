# velocity-cli

A Go component library for building styled CLI applications. Import it into any Go project for consistent output styling, tables, spinners, progress bars, and interactive prompts.

## Install

```bash
go get github.com/velocitykode/velocity-cli
```

## Output

```go
import cli "github.com/velocitykode/velocity-cli"

cli.Header("migrate")              // MIGRATE (styled header)
cli.Info("Running migrations...")   // → Running migrations...
cli.Success("Done")                 // ✓ Done
cli.Warning("No migrations found")  // ! No migrations found
cli.Error("Connection failed")      // ✗ Connection failed
cli.Muted("skipping...")            // dimmed text
cli.Bold("important")               // bold text

cli.Note("Server running on port 4000")  // boxed message (primary border)
cli.Alert("Database connection failed")  // boxed message (error border)
```

## Table

```go
cli.Table(
    []string{"Method", "Path", "Name"},
    [][]string{
        {"GET", "/users", "users.index"},
        {"POST", "/users", "users.store"},
    },
)
```

## Spinner & Progress

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

## Dependencies

- [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) — styling
- [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) — interactive prompts
- [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) — pre-built components

Zero framework dependency.

## License

MIT

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/banner"
	"github.com/velocitykode/velocity-cli/internal/colors"
)

var (
	// Use brand colors
	commandStyle = colors.SuccessStyle
	descStyle    = colors.MutedStyle
	sectionStyle = colors.BrandStyle.MarginTop(1)
)

// InitHelp sets up custom help for the root command
func InitHelp(root *cobra.Command) {
	root.SetHelpFunc(customHelpFunc)
	root.SetUsageFunc(customUsageFunc)
}

// Command groups for organized help output
var commandGroups = map[string][]string{
	"project":     {"new", "init"},
	"development": {"serve", "build"},
	"database":    {"migrate", "migrate:fresh"},
	"make":        {"make:controller"},
	"key":         {"key:generate"},
}

var groupOrder = []string{"project", "development", "database", "make", "key"}

func customHelpFunc(cmd *cobra.Command, args []string) {
	w := cmd.OutOrStdout()

	// Only show banner for root command
	if !cmd.HasParent() {
		fmt.Fprintln(w, banner.MediumBlocky())
		fmt.Fprintln(w)
	}

	// Usage
	fmt.Fprintln(w, sectionStyle.Render("Usage:"))
	fmt.Fprintf(w, "  %s [command]\n", cmd.CommandPath())
	fmt.Fprintln(w)

	// Commands - grouped by namespace for root command
	if len(cmd.Commands()) > 0 {
		if !cmd.HasParent() {
			// Root command - show grouped
			cmdMap := make(map[string]*cobra.Command)
			for _, c := range cmd.Commands() {
				if !c.Hidden {
					cmdMap[c.Name()] = c
				}
			}

			// Print groups
			for _, group := range groupOrder {
				cmds := commandGroups[group]
				hasCommands := false
				for _, name := range cmds {
					if _, ok := cmdMap[name]; ok {
						hasCommands = true
						break
					}
				}
				if !hasCommands {
					continue
				}

				fmt.Fprintln(w, descStyle.Render(group))
				for _, name := range cmds {
					if c, ok := cmdMap[name]; ok {
						fmt.Fprintf(w, "  %s  %s\n",
							commandStyle.Width(22).Render(c.Name()),
							descStyle.Render(c.Short))
						delete(cmdMap, name)
					}
				}
				fmt.Fprintln(w)
			}

			// Print remaining ungrouped commands (config, help, version)
			if len(cmdMap) > 0 {
				for _, c := range cmd.Commands() {
					if _, ok := cmdMap[c.Name()]; ok {
						fmt.Fprintf(w, "  %s  %s\n",
							commandStyle.Width(22).Render(c.Name()),
							descStyle.Render(c.Short))
					}
				}
				fmt.Fprintln(w)
			}
		} else {
			// Subcommand - show flat list
			fmt.Fprintln(w, sectionStyle.Render("Commands:"))
			for _, c := range cmd.Commands() {
				if !c.Hidden {
					fmt.Fprintf(w, "  %s  %s\n",
						commandStyle.Width(22).Render(c.Name()),
						descStyle.Render(c.Short))
				}
			}
			fmt.Fprintln(w)
		}
	}

	// Flags
	if cmd.HasAvailableLocalFlags() {
		fmt.Fprintln(w, sectionStyle.Render("Flags:"))
		fmt.Fprintln(w, descStyle.Render(cmd.LocalFlags().FlagUsages()))
	}

	// Footer
	fmt.Fprintln(w, descStyle.Render("Use \"velocity [command] --help\" for more information about a command"))
	fmt.Fprintln(w)
}

func customUsageFunc(cmd *cobra.Command) error {
	customHelpFunc(cmd, []string{})
	return nil
}

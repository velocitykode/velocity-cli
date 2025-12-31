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

	// Commands
	if len(cmd.Commands()) > 0 {
		fmt.Fprintln(w, sectionStyle.Render("Commands:"))
		for _, c := range cmd.Commands() {
			if !c.Hidden {
				fmt.Fprintf(w, "  %s  %s\n",
					commandStyle.Width(24).Render(c.Name()),
					descStyle.Render(c.Short))
			}
		}
		fmt.Fprintln(w)
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

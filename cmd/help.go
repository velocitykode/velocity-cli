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
	// Show banner
	fmt.Println(banner.MediumBlocky())
	fmt.Println()

	// Usage
	fmt.Println(sectionStyle.Render("Usage:"))
	fmt.Printf("  %s [command]\n", cmd.CommandPath())
	fmt.Println()

	// Commands
	if len(cmd.Commands()) > 0 {
		fmt.Println(sectionStyle.Render("Commands:"))
		for _, c := range cmd.Commands() {
			if !c.Hidden {
				fmt.Printf("  %s  %s\n",
					commandStyle.Width(24).Render(c.Name()),
					descStyle.Render(c.Short))
			}
		}
		fmt.Println()
	}

	// Flags
	if cmd.HasAvailableLocalFlags() {
		fmt.Println(sectionStyle.Render("Flags:"))
		fmt.Println(descStyle.Render(cmd.LocalFlags().FlagUsages()))
	}

	// Footer
	fmt.Println(descStyle.Render("Use \"velocity [command] --help\" for more information about a command"))
	fmt.Println()
}

func customUsageFunc(cmd *cobra.Command) error {
	customHelpFunc(cmd, []string{})
	return nil
}

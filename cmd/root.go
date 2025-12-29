package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/colors"
)

// Version is the CLI version - single source of truth
const Version = "0.2.0"

// Remove duplicate root command - it's defined in main.go

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		// Show compact banner for version
		fmt.Println(colors.BrandStyle.Render("╔══════════════════════════════════════╗"))
		fmt.Println(colors.BrandStyle.Render("║        VELOCITY CLI v" + Version + "        ║"))
		fmt.Println(colors.BrandStyle.Render("╚══════════════════════════════════════╝"))
		fmt.Println()
	},
}

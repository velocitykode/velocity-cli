package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/velocitykode/velocity-cli/internal/colors"
)

var version = "0.1.0"

// Remove duplicate root command - it's defined in main.go

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		// Show compact banner for version
		fmt.Println(colors.BrandStyle.Render("╔══════════════════════════════════════╗"))
		fmt.Println(colors.BrandStyle.Render("║        VELOCITY CLI v" + version + "        ║"))
		fmt.Println(colors.BrandStyle.Render("╚══════════════════════════════════════╝"))
		fmt.Println()
	},
}

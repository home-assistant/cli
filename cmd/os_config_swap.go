package cmd

import (
	"github.com/spf13/cobra"
)

var osConfigSwapCmd = &cobra.Command{
	Use:     "swap",
	Aliases: []string{"sw"},
	Short:   "Show or change Home Assistant OS swap settings",
	Long: `
This command allows you to show or change current swap configuration
of Home Assistant OS.`,
	Example: `
  ha os config swap info`,
}

func init() {
	osConfigCmd.AddCommand(osConfigSwapCmd)
}

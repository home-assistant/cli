package cmd

import (
	"github.com/spf13/cobra"
)

var osConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf", "cfg"},
	Short:   "Show or change Home Assistant OS settings",
	Long: `
This command allows you to show or change settings of Home Assistant OS.`,
	Example: `
  ha os config swap`,
}

func init() {
	osCmd.AddCommand(osConfigCmd)
}

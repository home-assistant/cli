package cmd

import (
	"github.com/spf13/cobra"
)

var hassosCmd = &cobra.Command{
	Use:     "hassos",
	Aliases: []string{"os"},
	Short:    "HassOS specific for updating, info and configuration imports",
	Long: `
This command set is specifically designed for HassOS and only works on those
systems. It provides an interface to get information about your HassOS system,
but also provides command to upgrade HassOS and the HassOS CLI. Finally,
it provides a command to import configurations from an USB-stick.`,
	Example: `
  hassio hassos info
  hassio hassos update`,
}

func init() {
	rootCmd.AddCommand(hassosCmd)
}

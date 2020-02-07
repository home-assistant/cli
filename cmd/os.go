package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var osCmd = &cobra.Command{
	Use:     "os",
	Aliases: []string{"hassos"},
	Short:   "Operating System specific for updating, info and configuration imports",
	Long: `
This command set is specifically designed for the Home Assistant Operating System
and only works on those systems. It provides an interface to get information
about your Home Assistant Operating System, but also provides command to
upgrade the operating system and the operating system CLI. Finally,
it provides a command to import configurations from an USB-stick.`,
	Example: `
  ha os info
	ha os update`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		for _, arg := range os.Args {
			if arg == "hassos" {
				cmd.PrintErrf("The use of '%s' is deprecated, please use 'os' instead!\n", arg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(osCmd)
}

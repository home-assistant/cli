package cmd

import (
	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"ho"},
	Short:   "Control the host/system that Home Assistant is running on",
	Long: `
The host command provides commandline tools to control the host (system) that
Home Assistant is running on. It allows you do thing like reboot or shutdown the
system, but also provides option to change the hostname of the system.`,
	Example: `
  ha host reboot
  ha host options --hostname "homeassistant.local"
`,
}

func init() {
	rootCmd.AddCommand(hostCmd)
}

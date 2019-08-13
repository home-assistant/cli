package cmd

import (
	"github.com/spf13/cobra"
)

var homeassistantCmd = &cobra.Command{
	Use:     "homeassistant",
	Aliases: []string{"home-assistant", "ha"},
	Short:   "Provides control of Home Assistant running on Hass.io",
	Long: `
This command provides a set of subcommands to control the Home Assistant
instance running on this Hass.io installation.

It provides commands to control Home Assistant (start, stop, restart), but also
allows you to check your Home Assistant configuration. Furthermore, some options
can be set and allows for upgrading/downgrading Home Assistant.
`,
	Example: `
  hassio homeassistant check
  hassio homeassistant restart
  hassio homeassistant update
  hassio homeassistant update --version 0.97.2`,
}

func init() {
	// add cmd to root command
	rootCmd.AddCommand(homeassistantCmd)
}

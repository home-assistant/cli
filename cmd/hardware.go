package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hardwareCmd = &cobra.Command{
	Use:     "hardware",
	Aliases: []string{"hw"},
	Short:   "Provides hardware information about your system",
	Long: `
The hardware command provides information about the hardware of your system
that is running Hass.io. It is useful for finding things like: available
audio devices and serial ports.`,
	Example: `
  hassio hardware info
  hassio hardware audio`,
}

func init() {
	log.Debug("Init hardware")

	// add cmd to root command
	rootCmd.AddCommand(hardwareCmd)
}

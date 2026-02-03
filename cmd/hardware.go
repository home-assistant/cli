package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var hardwareCmd = &cobra.Command{
	Use:     "hardware",
	Aliases: []string{"hw"},
	Short:   "Provides hardware information about your system",
	Long: `
The hardware command provides information about the hardware of your system
that is running Home Assistant. It is useful for finding things like: available
audio devices and serial ports.`,
	Example: `
  ha hardware info
  ha hardware audio`,
}

func init() {
	slog.Debug("Init hardware")

	// add cmd to root command
	rootCmd.AddCommand(hardwareCmd)
}

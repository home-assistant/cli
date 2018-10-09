package cmd

import (
	"github.com/spf13/cobra"
	"github.com/home-assistant/hassio-cli/cmd/hardware"
)

// hardwareCmd represents the hardware command
var hardwareCmd = &cobra.Command{
	Use:   "hardware",
	Aliases: []string{"ha"},
	Run: func(cmd *cobra.Command, args []string) {
		hardware.Execute()
	},
}

func init() {
	rootCmd.AddCommand(hardwareCmd)
}

package cmd

import (
	"github.com/spf13/cobra"
)

var osBoardsRaspberrypiCmd = &cobra.Command{
	Use:     "raspberrypi",
	Aliases: []string{"rpi"},
	Short:   "See or change settings of the current Raspberry Pi board",
	Long: `
This command allows you to see or change settings of the Raspberry Pi board that
Home Assistant is running on.`,
	Example: `
  ha os boards raspberrypi firmware`,
}

func init() {
	osBoardsCmd.AddCommand(osBoardsRaspberrypiCmd)
}

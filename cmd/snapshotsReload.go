package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// snapshotsReloadCmd represents the reload subcommand for snapshots
var snapshotsReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reload called")
	},
}

func init() {
	snapshotsCmd.AddCommand(snapshotsReloadCmd)
}

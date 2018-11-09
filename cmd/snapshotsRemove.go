package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// snapshotsRemoveCmd represents the remove subcommand for snapshots
var snapshotsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remove called")
	},
}

func init() {
	snapshotsCmd.AddCommand(snapshotsRemoveCmd)
}

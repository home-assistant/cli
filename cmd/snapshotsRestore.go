package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// snapshotsRestoreCmd represents the restore subcommand for snapshots
var snapshotsRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restore called")
	},
}

func init() {
	snapshotsCmd.AddCommand(snapshotsRestoreCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// snapshotsListCmd represents the list subcommand for snapshots
var snapshotsListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	snapshotsCmd.AddCommand(snapshotsListCmd)
}

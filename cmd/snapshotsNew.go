package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// snapshotsNewCmd represents the new subcommand for snapshots
var snapshotsNewCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
	},
}

func init() {
	snapshotsCmd.AddCommand(snapshotsNewCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// snapshotsInfoCmd represents the info subcommand for snapshots
var snapshotsInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("info called")
	},
}

func init() {
	snapshotsCmd.AddCommand(snapshotsInfoCmd)
}

package cmd

import (
	"github.com/home-assistant/hassio-cli/cmd/snapshots"
	"github.com/spf13/cobra"
)

// snapshotsCmd represents the snapshots command
var snapshotsCmd = &cobra.Command{
	Use:   "snapshots",
	Aliases: []string{"sa"},
	Run: func(cmd *cobra.Command, args []string) {
		snapshots.Execute()
	},
}

func init() {
	rootCmd.AddCommand(snapshotsCmd)
}

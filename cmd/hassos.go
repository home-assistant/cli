package cmd

import (
	"github.com/home-assistant/hassio-cli/cmd/hassos"
	"github.com/spf13/cobra"
)

// hassosCmd represents the hassos command
var hassosCmd = &cobra.Command{
	Use:   "hassos",
	Aliases: []string{"os"},
	Run: func(cmd *cobra.Command, args []string) {
		hassos.Execute()
	},
}

func init() {
	rootCmd.AddCommand(hassosCmd)
}

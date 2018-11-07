package cmd

import (
	"github.com/home-assistant/hassio-cli/cmd/host"
	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"ho"},
	Run: func(cmd *cobra.Command, args []string) {
		host.Execute()
	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}

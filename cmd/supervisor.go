package cmd

import (
	"github.com/home-assistant/hassio-cli/cmd/supervisor"
	"github.com/spf13/cobra"
)

// supervisorCmd represents the supervisor command
var supervisorCmd = &cobra.Command{
	Use:   "supervisor",
	Aliases: []string{"su"},
	Run: func(cmd *cobra.Command, args []string) {
		supervisor.Execute()
	},
}

func init() {
	rootCmd.AddCommand(supervisorCmd)
}

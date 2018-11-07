package cmd

import (
	"github.com/home-assistant/hassio-cli/cmd/addons"
	"github.com/spf13/cobra"
)

// addonsCmd represents the addons command
var addonsCmd = &cobra.Command{
	Use:     "addons",
	Aliases: []string{"ad"},
	Run: func(cmd *cobra.Command, args []string) {
		addons.Execute()
	},
}

func init() {
	rootCmd.AddCommand(addonsCmd)

}

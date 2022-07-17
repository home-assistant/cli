package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var availableUpdatesCmd = &cobra.Command{
	Use:     "available-updates",
	Aliases: []string{"updates", "available_updates"},
	Short:   "Provides information about current pending updates",
	Long: `
This command provides information about the currently pending updates on the system.
	`,
	Example: `
  ha available-updates
	`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("available_updates")

		section := "available_updates"
		command := ""

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	rootCmd.AddCommand(availableUpdatesCmd)
}

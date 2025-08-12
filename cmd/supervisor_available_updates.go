package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var supervisorAvailableUpdatesCmd = &cobra.Command{
	Use:     "available-updates",
	Aliases: []string{"available_updates", "updates"},
	Short:   "Provides information about available updates",
	Long: `
This command provides you information about available updates.`,
	Example: `
  ha supervisor available-updates`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor available-updates")

		section := "supervisor"
		command := "available_updates"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}

	},
}

func init() {
	supervisorCmd.AddCommand(supervisorAvailableUpdatesCmd)
}

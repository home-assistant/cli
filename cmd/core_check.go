package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coreCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"validate", "chk", "ch"},
	Short:   "Validates your Home Assistant Core configuration",
	Long: `
This commands allows you to check/validate your, currently on disk stored,
Home Assistant Core configuration. This is helpful when you've made changes and
want to make sure the configuration is right, before restarting
Home Assistant Core.`,
	Example: `
  ha core check`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core check")

		section := "core"
		command := "check"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, nil, helper.ContainerOperationTimeout)
		ProgressSpinner.Stop()

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	coreCmd.AddCommand(coreCheckCmd)
}

package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var jobsResetCmd = &cobra.Command{
	Use:               "reset",
	Short:             "Resets the internal Home Assistant Job Manager configuration",
	Long:              `Resets the internal Home Assistant Job Manager configuration.`,
	Example:           `ha jobs reset`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("jobs reset", "args", args)

		section := "jobs"
		command := "reset"

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(section, command, nil)
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
	jobsCmd.AddCommand(jobsResetCmd)
}

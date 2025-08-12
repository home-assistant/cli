package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var resolutionCheckRunCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"execute", "ru"},
	Short:   "Run a specific check at the Resolution center",
	Long: `
This command executes an backend check immediately on the system.`,
	Example: `
  ha resolution check run [slug]`,
	ValidArgsFunction: resolutionCheckCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("check run")

		section := "resolution"
		command := "check/{check}/run"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		check := args[0]

		request.SetPathParams(map[string]string{
			"check": check,
		})

		ProgressSpinner.Start()
		resp, err := request.Post(url)
		ProgressSpinner.Stop()

		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	resolutionCheckCmd.AddCommand(resolutionCheckRunCmd)
}

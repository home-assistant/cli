package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var appsStatsCmd = &cobra.Command{
	Use:     "stats [slug]",
	Aliases: []string{"status", "stat"},
	Short:   "Provides system usage stats of a Home Assistant app",
	Long: `
Provides insight into the system usage stats of an app. It shows you
how much CPU, memory, disk & network resources it uses.
`,
	Example: `
  ha apps stats core_ssh
`,
	ValidArgsFunction: appsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("apps stats")

		section := "addons"
		command := "{slug}/stats"

		url, err := helper.URLHelper(section, command)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.Get(url)
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
	appsCmd.AddCommand(appsStatsCmd)
}

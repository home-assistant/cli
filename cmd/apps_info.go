package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var appsInfoCmd = &cobra.Command{
	Use:     "info [slug]",
	Aliases: []string{"in", "info"},
	Short:   "Show information about available Home Assistant apps",
	Long: `
This command can provide information on all available apps or, if a slug
is provided, information about a specific app.
`,
	Example: `
  ha apps info
  ha apps info core_ssh
`,
	ValidArgsFunction: appsCompletions,
	Args:              cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("apps info")

		section := "addons"
		command := "{slug}/info"

		url, err := helper.URLHelper(section, command)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		slug := "self"
		if len(args) > 0 {
			slug = args[0]
		}

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
	appsCmd.AddCommand(appsInfoCmd)
}

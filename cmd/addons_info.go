package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsInfoCmd = &cobra.Command{
	Use:     "info [slug]",
	Aliases: []string{"in", "info"},
	Short:   "Show information about available Home Assistant add-ons",
	Long: `
This command can provide information on all available add-ons or, if a slug
is provided, information about a specific add-on.
`,
	Example: `
  ha addons info
  ha addons info core_ssh
`,
	ValidArgsFunction: addonsCompletions,
	Args:              cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons info")

		section := "addons"
		command := "{slug}/info"

		url, err := helper.URLHelper(section, command)

		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	addonsCmd.AddCommand(addonsInfoCmd)
}

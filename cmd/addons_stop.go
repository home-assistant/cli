package cmd

import (
	"fmt"

	"github.com/home-assistant/cli/client"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsStopCmd = &cobra.Command{
	Use:     "stop [slug]",
	Aliases: []string{"halt", "shutdown", "quit"},
	Short:   "Manually stop a running Home Assistant add-on",
	Long: `
This command allows you to manually start a stopped Home Assistant add-on
`,
	Example: `
  ha addons stop core_ssh
`,
	ValidArgsFunction: addonsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons stop")

		section := "addons"
		command := "{slug}/stop"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequestTimeout(helper.ContainerOperationTimeout)

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.Post(url)
		resp, err = client.GenericJSONErrorHandling(resp, err)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {

	addonsCmd.AddCommand(addonsStopCmd)
}

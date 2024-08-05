package cmd

import (
	"fmt"

	"github.com/home-assistant/cli/client"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsRestartCmd = &cobra.Command{
	Use:     "restart [slug]",
	Aliases: []string{"reboot"},
	Short:   "Restarts a Home Assistant add-on",
	Long: `
Restart a Home Assistant add-on
`,
	Example: `
  ha addons restart core_ssh
`,
	ValidArgsFunction: addonsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons restart")

		section := "addons"
		command := "{slug}/restart"

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

	addonsCmd.AddCommand(addonsRestartCmd)
}

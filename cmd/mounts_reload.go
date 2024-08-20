package cmd

import (
	"fmt"

	"github.com/home-assistant/cli/client"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var mountsReloadCmd = &cobra.Command{
	Use:     "reload [name]",
	Aliases: []string{"re", "refresh", "remount"},
	Short:   "Reload a mount in Supervisor",
	Long: `
Unmount and remount an existing mount in Supervisor using
the same configuration.
`,
	Example: `
  ha mounts reload my_share
`,
	ValidArgsFunction: mountsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("mounts reload")

		section := "mounts"
		command := "{name}/reload"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		name := args[0]
		request.SetPathParams(map[string]string{
			"name": name,
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
	mountsCmd.AddCommand(mountsReloadCmd)
}

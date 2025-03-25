package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var mountsUpdateCmd = &cobra.Command{
	Use:     "update [name]",
	Aliases: []string{"change", "set", "up", "modify", "mod"},
	Short:   "Update configuration of a mount in Supervisor",
	Long: `
Update or change the configuration of an existing mount in Supervisor.
`,
	Example: `
  ha mounts update my_share --usage media --type cifs --server server.local --share media
`,
	ValidArgsFunction: mountsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("mounts update")

		section := "mounts"
		command := "{name}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		name := args[0]
		options := make(map[string]any)

		request.SetPathParams(map[string]string{
			"name": name,
		})
		mountFlagsToOptions(cmd, options)

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		resp, err := request.Put(url)
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
	addMountFlags(mountsUpdateCmd)
	mountsCmd.AddCommand(mountsUpdateCmd)
}

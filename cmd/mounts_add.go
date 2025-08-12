package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var mountsAddCmd = &cobra.Command{
	Use:     "add [name]",
	Aliases: []string{"create", "new"},
	Short:   "Add new mount to Supervisor",
	Long: `
Add and configure a new mount in Supervisor.
`,
	Example: `
  ha mounts add my_share --usage media --type cifs --server server.local --share media
`,
	ValidArgsFunction: mountsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("mounts add")

		section := "mounts"
		command := ""

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		name := args[0]
		options := make(map[string]any)

		options["name"] = name
		mountFlagsToOptions(cmd, options)

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		resp, err := request.Post(url)
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
	addMountFlags(mountsAddCmd)
	mountsCmd.AddCommand(mountsAddCmd)
}

package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
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
  ha mounts add media_disk --usage media --type disk --device /dev/sda1
`,
	ValidArgsFunction: mountsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("mounts add", "args", args)

		section := "mounts"
		command := ""

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequestTimeout(helper.MountTimeout)

		name := args[0]
		options := make(map[string]any)

		options["name"] = name
		mountFlagsToNewOptions(cmd, options)

		if len(options) > 0 {
			slog.Debug("Request body", "options", options)
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

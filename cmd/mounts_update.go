package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
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
		slog.Debug("mounts update", "args", args)

		section := "mounts"
		command := "{name}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
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
			slog.Debug("Request body", "options", options)
			request.SetBody(options)
		}

		resp, err := request.Put(url)
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
	addMountFlags(mountsUpdateCmd)
	mountsCmd.AddCommand(mountsUpdateCmd)
}

package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var mountsDeleteCmd = &cobra.Command{
	Use:     "delete [name]",
	Aliases: []string{"del", "remove", "rm"},
	Short:   "Delete a mount from Supervisor",
	Long: `
Unmount and delete an existing mount from Supervisor.
`,
	Example: `
  ha mounts delete my_share
`,
	ValidArgsFunction: mountsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("mounts delete", "args", args)

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
		request.SetPathParams(map[string]string{
			"name": name,
		})

		resp, err := request.Delete(url)
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
	mountsCmd.AddCommand(mountsDeleteCmd)
}

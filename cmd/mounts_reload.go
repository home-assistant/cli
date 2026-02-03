package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
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
		slog.Debug("mounts reload", "args", args)

		section := "mounts"
		command := "{name}/reload"

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
	mountsCmd.AddCommand(mountsReloadCmd)
}

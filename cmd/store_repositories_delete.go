package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var storeRepositoriesDeleteCmd = &cobra.Command{
	Use:     "delete [slug]",
	Aliases: []string{"del", "remove"},
	Short:   "Delete repository from Home Assistant store",
	Long: `
Remove a repository of apps that isn't in use from the Home Assistant store.
`,
	Example: `
ha store delete 94cfad5a
`,
	ValidArgsFunction: storeRepositoriesCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("store delete", "args", args)

		section := "store"
		command := "repositories/{slug}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		slug := args[0]
		request.SetPathParams(map[string]string{
			"slug": slug,
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
	storeCmd.AddCommand(storeRepositoriesDeleteCmd)
}

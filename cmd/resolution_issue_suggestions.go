package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var resolutionIssueSuggestionsCmd = &cobra.Command{
	Use:     "suggestions",
	Aliases: []string{"su", "solutions"},
	Short:   "Suggestions which resolve an issue",
	Long: `
This command returns suggestions which resolve an issue when applied.`,
	Example: `
  ha resolution issue suggestions [id]`,
	ValidArgsFunction: resolutionIssueCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("issue suggestions", "args", args)

		section := "resolution"
		command := "issue/{issue}/suggestions"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		issue := args[0]

		request.SetPathParams(map[string]string{
			"issue": issue,
		})

		resp, err := request.Get(url)
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
	resolutionIssueCmd.AddCommand(resolutionIssueSuggestionsCmd)
}

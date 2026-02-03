package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var resolutionSuggestionDismissCmd = &cobra.Command{
	Use:     "dismiss",
	Aliases: []string{"disable", "remove"},
	Short:   "Suggestion dismiss reported by Resolution center",
	Long: `
This command allows dismissing a suggestion reported by the system.`,
	Example: `
  ha resolution suggestion dismiss [id]`,
	ValidArgsFunction: resolutionSuggestionCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("suggestion dismiss", "args", args)

		section := "resolution"
		command := "suggestion/{suggestion}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		suggestion := args[0]

		request.SetPathParams(map[string]string{
			"suggestion": suggestion,
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
	resolutionSuggestionCmd.AddCommand(resolutionSuggestionDismissCmd)
}

package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var resolutionSuggestionApplyCmd = &cobra.Command{
	Use:     "apply",
	Aliases: []string{"enable", "run"},
	Short:   "Suggestion apply reported by Resolution center",
	Long: `
This command allow to apply an suggestion reported by the System.`,
	Example: `
  ha resolution suggestion apply [id]`,
	ValidArgsFunction: resolutionSuggestionCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("suggestion dismiss")

		section := "resolution"
		command := "suggestion/{suggestion}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		suggestion := args[0]

		request.SetPathParams(map[string]string{
			"suggestion": suggestion,
		})

		resp, err := request.Post(url)
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
	resolutionSuggestionCmd.AddCommand(resolutionSuggestionApplyCmd)
}

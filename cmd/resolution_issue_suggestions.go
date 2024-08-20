package cmd

import (
	"fmt"

	"github.com/home-assistant/cli/client"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
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
		log.WithField("args", args).Debug("issue suggestions")

		section := "resolution"
		command := "issue/{issue}/suggestions"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		issue := args[0]

		request.SetPathParams(map[string]string{
			"issue": issue,
		})

		resp, err := request.Get(url)
		resp, err = client.GenericJSONErrorHandling(resp, err)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	resolutionIssueCmd.AddCommand(resolutionIssueSuggestionsCmd)
}

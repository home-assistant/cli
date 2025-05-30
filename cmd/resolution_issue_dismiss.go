package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var resolutionIssueDismissCmd = &cobra.Command{
	Use:     "dismiss",
	Aliases: []string{"disable", "remove"},
	Short:   "Dismiss issues",
	Long: `
This command allows dismissing issues reported by the system.`,
	Example: `
  ha resolution issue dismiss [id]`,
	ValidArgsFunction: resolutionIssueCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("issue dismiss")

		section := "resolution"
		command := "issue/{issue}"

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

		resp, err := request.Delete(url)
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
	resolutionIssueCmd.AddCommand(resolutionIssueDismissCmd)
}

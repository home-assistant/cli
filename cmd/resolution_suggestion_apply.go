package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
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
	ValidArgsFunction: cobra.NoFileCompletions,
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

		// returns 200 OK or 400, everything else is wrong
		if err == nil {
			if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
				err = errors.New("Unexpected server response")
				log.Error(err)
			} else if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
				err = errors.New("API did not return a JSON response")
				log.Error(err)
			}
		}

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

package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resolutionSuggestionDismissCmd = &cobra.Command{
	Use:     "dismiss",
	Aliases: []string{"disable", "remove"},
	Short:   "Suggestion dismiss reported by Resolution center",
	Long: `
This command allows dismissing a suggestion reported by the system.`,
	Example: `
  ha resolution suggestion dismiss [id]`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("suggestion dismiss")

		section := "resolution"
		command := "suggestion/{suggestion}"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
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

		resp, err := request.Delete(url)

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
	resolutionSuggestionCmd.AddCommand(resolutionSuggestionDismissCmd)
}

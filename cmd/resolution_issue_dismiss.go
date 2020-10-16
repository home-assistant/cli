package cmd

import (
	"errors"
	"fmt"
	"net/http"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resolutionIssueDismissCmd = &cobra.Command{
	Use:     "dismiss",
	Aliases: []string{"disable", "remove"},
	Short:   "Dismiss issues",
	Long: `
This command allows dismissing issues reported by the system.`,
	Example: `
  ha resolution issue dismiss [id]`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("issue dismiss")

		section := "resolution"
		command := "issue/{issue}"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
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

		// returns 200 OK or 400, everything else is wrong
		if err == nil {
			if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
				err = errors.New("Unexpected server response")
				log.Error(err)
			} else if !resty.IsJSONType(resp.Header().Get(http.CanonicalHeaderKey("Content-Type"))) {
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

		return
	},
}

func init() {
	resolutionIssueCmd.AddCommand(resolutionIssueDismissCmd)
}

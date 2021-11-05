package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsStatsCmd = &cobra.Command{
	Use:     "stats [slug]",
	Aliases: []string{"status", "stat"},
	Short:   "Provides system usage stats of a Home Assistant add-on",
	Long: `
Provides insight into the system usage stats of an add-on. It shows you
how much CPU, memory, disk & network resources it uses.
`,
	Example: `
  ha addons stats core_ssh
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons stats")

		section := "addons"
		command := "{slug}/stats"

		url, err := helper.URLHelper(section, command)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.Get(url)

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
	addonsCmd.AddCommand(addonsStatsCmd)
}

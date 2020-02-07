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

var addonsStartCmd = &cobra.Command{
	Use:     "start [slug]",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"run", "st"},
	Short:   "Manually start a stopped Home Assistant add-on",
	Long: `
This command allows you to manually start a stopped Home Assistant add-on
`,
	Example: `
  ha addons start core_ssh
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons start")

		section := "addons"
		command := "{slug}/start"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
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

		resp, err := request.Post(url)

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

	addonsCmd.AddCommand(addonsStartCmd)
}

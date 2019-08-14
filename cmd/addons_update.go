package cmd

import (
	"errors"
	"fmt"
	"net/http"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	resty "github.com/go-resty/resty/v2"
)

var addonsUpdateCmd = &cobra.Command{
	Use:  "update [slug]",
	Args: cobra.ExactArgs(1),
	Aliases: []string{"upgrade", "up"},
	Short:   "Upgrades an Hass.io add-on to the latest version",
	Long: `
Using this command you can upgrade an Hass.io add-on to its latest version.
It is currently not possible to upgrade/downgrade to a specific version.
`,
	Example: `
  hassio addons update core_ssh
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons update")

		section := "addons"
		command := "{slug}/update"
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

	addonsCmd.AddCommand(addonsUpdateCmd)
}

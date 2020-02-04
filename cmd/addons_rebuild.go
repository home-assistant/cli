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

var addonsRebuildCmd = &cobra.Command{
	Use:     "rebuild [slug]",
	Aliases: []string{"rb", "reinstall"},
	Short:   "Rebuild a locally build Home Assistant add-on",
	Long: `
Most add-ons provide pre-build images Home Assistant can download an use.
However, some don't. This is usually the case for local or development version
of add-ons. This command allows you to trigger a rebuild of a locally build
add-on.
`,
	Example: `
  ha addons rebuild local_my_addon
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons rebuild")

		section := "addons"
		command := "{slug}/rebuild"
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

		ProgressSpinner.Start()
		resp, err := request.Post(url)
		ProgressSpinner.Stop()

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

	addonsCmd.AddCommand(addonsRebuildCmd)
}

package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeAddonsInstalCmd = &cobra.Command{
	Use:     "install [slug]",
	Aliases: []string{"i", "inst"},
	Short:   "Installs a Home Assistant add-on",
	Long: `
This command allows you to install a Home Assistant add-on from the commandline.
`,
	Example: `
  ha store addons install core_ssh
`,
	ValidArgsFunction: storeAddonCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store addons install")

		section := "store"
		command := "addons/{slug}/install"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequestTimeout(helper.ContainerDownloadTimeout)

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

	storeAddonsCmd.AddCommand(storeAddonsInstalCmd)
}

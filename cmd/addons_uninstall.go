package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsUninstallCmd = &cobra.Command{
	Use:     "uninstall [slug]",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"remove", "delete", "del", "rem", "un", "uninst"},
	Short:   "Uninstalls a Home Assistant add-on",
	Long: `
This command allows you to uninstall a Home Assistant add-on.
`,
	Example: `
  ha addons uninstall core_ssh
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons uninstall")

		section := "addons"
		command := "{slug}/uninstall"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequestTimeout(helper.ContainerOperationTimeout)

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

	addonsCmd.AddCommand(addonsUninstallCmd)
}

package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeAddonsUpdateCmd = &cobra.Command{
	Use:     "update [slug]",
	Aliases: []string{"upgrade", "up"},
	Short:   "Upgrades a Home Assistant add-on to the latest version",
	Long: `
This command can upgrade a Home Assistant add-on to its latest version.
It is currently not possible to upgrade/downgrade to a specific version.
`,
	Example: `
  ha store addons update core_ssh
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store addons update")

		section := "store"
		command := "addons/{slug}/update"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		ProgressSpinner.Start()
		request := helper.GetJSONRequest()
		ProgressSpinner.Stop()

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		options := make(map[string]interface{})

		backup, _ := cmd.Flags().GetBool("backup")
		if cmd.Flags().Changed("backup") {
			request.SetBody(options)
			options["backup"] = backup
		}

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

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
	storeAddonsUpdateCmd.Flags().Bool("backup", false, "Create partial backup before update")
	storeAddonsUpdateCmd.RegisterFlagCompletionFunc("backup", boolCompletions)
	storeAddonsCmd.AddCommand(storeAddonsUpdateCmd)
}

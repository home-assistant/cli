package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeRepositoriesAddCmd = &cobra.Command{
	Use:     "add [repository]",
	Aliases: []string{"set", "new"},
	Short:   "Add new repository to Home Assistant store",
	Long: `
Add new repository of add-ons to the Home Assistant store.
`,
	Example: `
ha store add https://github.com/home-assistant/addons-example
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store add")

		section := "store"
		command := "repositories"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		options := map[string]string{}

		request := helper.GetJSONRequest()

		repository := args[0]
		options["repository"] = repository

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		resp, err := request.Post(url)

		// returns 200 OK or 400, everything else is wrong
		if err == nil {
			if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
				err = errors.New("unexpected server response")
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
	storeCmd.AddCommand(storeRepositoriesAddCmd)
}

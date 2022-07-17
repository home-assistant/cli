package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeRepositoriesDeleteCmd = &cobra.Command{
	Use:     "delete [slug]",
	Aliases: []string{"del", "remove"},
	Short:   "Delete repository from Home Assistant store",
	Long: `
Remove a repository of add-ons that isn't in use from the Home Assistant store.
`,
	Example: `
ha repositories delete 94cfad5a
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store delete")

		section := "store"
		command := "repositories/{slug}"

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

		resp, err := request.Delete(url)

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
	storeCmd.AddCommand(storeRepositoriesDeleteCmd)
}

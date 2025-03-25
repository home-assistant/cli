package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var authCacheCmd = &cobra.Command{
	Use:     "cache",
	Aliases: []string{"data", "ca"},
	Short:   "Reset the auth cache of Home Assistant on Supervisor.",
	Long: `
This command allows you to reset the internal password cache of a Home Assistant auth.
`,
	Example: `
  ha authentication cache
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("auth cache")

		section := "auth"
		command := "cache"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		resp, err := request.Delete(url)
		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	authCmd.AddCommand(authCacheCmd)
}

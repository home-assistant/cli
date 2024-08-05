package cmd

import (
	"fmt"

	"github.com/home-assistant/cli/client"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addonsUninstallCmd = &cobra.Command{
	Use:     "uninstall [slug]",
	Aliases: []string{"remove", "delete", "del", "rem", "un", "uninst"},
	Short:   "Uninstalls a Home Assistant add-on",
	Long: `
This command allows you to uninstall a Home Assistant add-on.
`,
	Example: `
  ha addons uninstall core_ssh
`,
	ValidArgsFunction: addonsCompletions,
	Args:              cobra.ExactArgs(1),
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

		var options map[string]interface{}

		removeConfig, _ := cmd.Flags().GetBool("remove-config")
		if removeConfig {
			options = map[string]interface{}{"remove_config": removeConfig}
			request.SetBody(options)
		}

		resp, err := request.Post(url)
		resp, err = client.GenericJSONErrorHandling(resp, err)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	addonsUninstallCmd.Flags().Bool("remove-config", false, "Delete addon's config folder (if used)")
	addonsUninstallCmd.Flags().Lookup("remove-config").NoOptDefVal = "true"
	addonsUninstallCmd.RegisterFlagCompletionFunc("remove-config", boolCompletions)

	addonsCmd.AddCommand(addonsUninstallCmd)
}

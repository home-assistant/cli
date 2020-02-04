package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var osUpdateCliCmd = &cobra.Command{
	Use:     "update-cli",
	Aliases: []string{"up-cli", "upcli", "cli-update", "cli-up", "cliup"},
	Short:   "Updates the Home Assistant Operating System CLI",
	Long: `
Using this command you can upgrade or downgrade the Home Assistant
Operating System CLI to the latest version or the version specified.
`,
	Example: `
  ha os update-cli --version 5
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os update-cli")

		section := "hassos"
		command := "update/cli"
		base := viper.GetString("endpoint")

		var options map[string]interface{}

		version, err := cmd.Flags().GetString("version")
		if version != "" {
			options = map[string]interface{}{"version": version}
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(base, section, command, options)
		ProgressSpinner.Stop()
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
	osUpdateCliCmd.Flags().StringP("version", "", "", "Version to update to")
	osCmd.AddCommand(osUpdateCliCmd)
}

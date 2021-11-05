package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var observerUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade", "downgrade", "up", "down"},
	Short:   "Updates the internal Home Assistant observer",
	Long: `
Using this command you can upgrade or downgrade the internal Home Assistant 
observer, to the latest version or the version specified.
`,
	Example: `
  ha observer update --version 5
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("observer update")

		section := "observer"
		command := "update"

		var options map[string]interface{}

		version, _ := cmd.Flags().GetString("version")
		if version != "" {
			options = map[string]interface{}{"version": version}
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, options, helper.ContainerDownloadTimeout)
		ProgressSpinner.Stop()

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	observerUpdateCmd.Flags().StringP("version", "", "", "Version to update to")
	observerCmd.AddCommand(observerUpdateCmd)
}

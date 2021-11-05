package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade", "downgrade", "up", "down"},
	Short:   "Updates the Home Assistant Operating System",
	Long: `
Using this command you can upgrade or downgrade the Home Assistant 
Operating System to the latest version or the version specified.
`,
	Example: `
  ha os update
  ha os update --version 5
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os update")

		section := "os"
		command := "update"

		var options map[string]interface{}

		version, _ := cmd.Flags().GetString("version")
		if version != "" {
			options = map[string]interface{}{"version": version}
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, options, helper.OsDownloadTimeout)
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
	osUpdateCmd.Flags().StringP("version", "", "", "Version to update to")
	osCmd.AddCommand(osUpdateCmd)
}

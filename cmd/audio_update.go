package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var audioUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade", "downgrade", "up", "down"},
	Short:   "Updates the Home Assistant Audio",
	Long: `
Using this command you can upgrade or downgrade the Home Assistant Audio
instance running on your system to the latest version or the version specified.`,
	Example: `
  ha audio update
  ha audio update --version 6`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("audio update")

		section := "audio"
		command := "update"
		base := viper.GetString("endpoint")

		var options map[string]interface{}

		version, _ := cmd.Flags().GetString("version")
		if version != "" {
			options = map[string]interface{}{"version": version}
		}

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(base, section, command, options, helper.ContainerDownloadTimeout)
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
	audioUpdateCmd.Flags().StringP("version", "", "", "Version to update to")
	audioCmd.AddCommand(audioUpdateCmd)
}

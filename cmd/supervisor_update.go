package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supervisorUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade", "downgrade", "up", "down"},
	Short:   "Updates the Home Assistant Supervisor",
	Long: `
Using this command you can upgrade or downgrade the Home Assistant Supervisor
running on your Home Assistant  system to the latest version
or the version specified.`,
	Example: `
  ha supervisor update
  ha supervisor update --version 173`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor update")

		section := "supervisor"
		command := "update"
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
	supervisorUpdateCmd.Flags().StringP("version", "", "", "Version to update to")
	supervisorCmd.AddCommand(supervisorUpdateCmd)
}

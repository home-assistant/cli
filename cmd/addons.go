package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addonsCmd = &cobra.Command{
	Use:     "addons",
	Aliases: []string{"addon", "add-on", "add-ons", "ad"},
	Short:   "Install, update, remove and configure Home Assistant add-ons",
	Long: `
The addons command allows you to manage Home Assistant add-ons by exposing
commands for installing, removing, configure and control them. It also provides
information commands for add-ons.`,
	Example: `
  ha addons logs core_ssh
  ha addons install core_ssh
  ha addons start core_ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("addons")

		section := "addons"
		command := ""
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
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
	log.Debug("Init addons")

	rootCmd.AddCommand(addonsCmd)
}

package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeAddonsCmd = &cobra.Command{
	Use:     "addons",
	Aliases: []string{"add-on", "addon", "add-ons"},
	Short:   "Install and update Home Asistant add-ons",
	Long: `
The store command allows you to manage Home Assistant add-ons by exposing
commands for installing or update them.`,
	Example: `
  ha store addons install core_ssh
  ha store addons update core_ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("store")

		section := "store"
		command := ""

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	storeCmd.AddCommand(storeAddonsCmd)
}

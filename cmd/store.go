package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Use:     "store",
	Aliases: []string{"shop", "stor"},
	Short:   "Install and update Home Assistant add-ons and manage stores",
	Long: `
The store command allows you to manage Home Assistant add-ons by exposing
commands for installing or update them. It also provides functionality
for managing stores that provide later add-ons.`,
	Example: `
  ha store addons install core_ssh
  ha store add https://github.com/home-assistant/addons-example
  ha store delete 94cfad5a 
  ha store reload`,
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
	log.Debug("Init store")

	rootCmd.AddCommand(storeCmd)
}

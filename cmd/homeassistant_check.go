package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeassistantCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"validate", "chk", "ch"},
	Short:   "Validates your Home Assistant configuration",
	Long: `
This commands allows you to check/validate your, currently on disk stored,
Home Assistant configuration. This is helpful when you've made changes and
want to make sure the configuration is right, before restarting
Home Assistant.`,
	Example: `
  hassio homeassistant check`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant check")

		section := "homeassistant"
		command := "check"
		base := viper.GetString("endpoint")
		
		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPost(base, section, command, nil)
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
	homeassistantCmd.AddCommand(homeassistantCheckCmd)
}

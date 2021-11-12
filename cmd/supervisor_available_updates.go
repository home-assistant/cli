package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supervisorAvailableUpdatesCmd = &cobra.Command{
	Use:     "available_updates",
	Aliases: []string{"updates"},
	Short:   "Provides information about available updates",
	Long: `
This command provides you information about available updates.`,
	Example: `
  ha supervisor info`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor available_updates")

		section := "supervisor"
		command := "available_updates"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}

	},
}

func init() {
	supervisorCmd.AddCommand(supervisorAvailableUpdatesCmd)
}

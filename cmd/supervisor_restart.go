package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supervisorRestartCmd = &cobra.Command{
	Use:     "restart",
	Aliases: []string{"reboot"},
	Short:   "Restart the Home Assistant Supervisor",
	Long: `
Restart the Supervisor internal, this can solve healthy issues.`,
	Example: `
  ha supervisor restart`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor restart")

		section := "supervisor"
		command := "restart"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	supervisorCmd.AddCommand(supervisorRestartCmd)
}

package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supervisorRepairCmd = &cobra.Command{
	Use:     "repair",
	Aliases: []string{"rep", "fix"},
	Short:   "Repair Docker issue automatically using the Supervisor (BETA!)",
	Long: `
There are cases where the Docker file system running on your Hass.io system,
encounters issue or corruptions. Running this command, the Hass.io Supervisor
will try to resolve these.
WARNING! This command is currently in beta.`,
	Example: `
  hassio supervisor repair`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor repair")

		section := "supervisor"
		command := "repair"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
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
	supervisorCmd.AddCommand(supervisorRepairCmd)
}

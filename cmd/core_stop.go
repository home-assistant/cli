package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var coreStopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{},
	Short:   "Manually stop Home Assistant Core",
	Long: `
This command allows you to manually stop the Home Assistant Core instance on
your system.`,
	Example: `
  ha core stop`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core stop")

		section := "homeassistant"
		command := "stop"
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
	coreCmd.AddCommand(coreStopCmd)
}

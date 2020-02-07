package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hardwareInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in", "inf"},
	Short:   "Provides hardware information about your system",
	Long: `
The hardware command provides information about the hardware of your system
that is running Home Assistant. It is useful for finding things like: available 
serial ports.`,
	Example: `ha hardware info`,

	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("hardware info")

		section := "hardware"
		command := "info"
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
	hardwareCmd.AddCommand(hardwareInfoCmd)
}

package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var hardwareAudioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"sounds", "snd", "au"},
	Short:   "Provides information about audio devices on your system",
	Long: `
The command provides information about audio devices available on your system.`,
	Example: `ha hardware audio`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("hardware info")

		section := "hardware"
		command := "audio"
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
	hardwareCmd.AddCommand(hardwareAudioCmd)
}

package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addonsCmd = &cobra.Command{
	Use:     "addons",
	Aliases: []string{"ad"},
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

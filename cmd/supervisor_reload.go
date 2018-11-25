package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var supervisorReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"re"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("supervisor reload")

		section := "supervisor"
		command := "reload"
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONPost(base, section, command, nil)
		if err != nil {
			fmt.Println(err)
		} else {
			helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	supervisorCmd.AddCommand(supervisorReloadCmd)
}

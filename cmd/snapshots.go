package cmd

import (
	"fmt"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var snapshotsCmd = &cobra.Command{
	Use:     "snapshots",
	Aliases: []string{"sn"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("snapshots")

		section := "snapshots"
		command := ""
		base := viper.GetString("endpoint")

		resp, err := helper.GenericJSONGet(base, section, command)
		if err != nil {
			fmt.Println(err)
		} else {
			helper.ShowJSONResponse(resp)
		}
		return
	},
}

func init() {
	log.Debug("Init snapshots")
	// add subcommands
	// TODO: add subcommand

	// add cmd to root command
	rootCmd.AddCommand(snapshotsCmd)
}

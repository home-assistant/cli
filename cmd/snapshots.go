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
	Aliases: []string{"snapshot", "snap", "shot", "sn", "backup", "backups", "bk"},
	Short:   "Create, restore and remove snapshot backups",
	Long: `
Snapshots are backups of your Hass.io system, which you can create, restore,
and delete using this command.`,
	Example: `
  hassio snapshots
  hassio snapshots new`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("snapshots")

		section := "snapshots"
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
	log.Debug("Init snapshots")
	// add cmd to root command
	rootCmd.AddCommand(snapshotsCmd)
}

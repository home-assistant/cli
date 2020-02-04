package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var snapshotsReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"refresh", "re"},
	Short:   "Reload the files on disk to check for new or removed snapshots",
	Long: `
If a snapshot has been manually placed inside the backup folder, or has been
removed manually, this command can trigger Home Assistant to re-read the files
on disk`,
	Example: `
  ha snapshots reload`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("snapshots reload")

		section := "snapshots"
		command := "reload"
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
	snapshotsCmd.AddCommand(snapshotsReloadCmd)
}

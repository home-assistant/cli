package cmd

import (
	"fmt"
	"os"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupsCmd = &cobra.Command{
	Use:     "backups",
	Aliases: []string{"snapshot", "snap", "shot", "sn", "backup", "backups", "bk"},
	Short:   "Create, restore and remove backups",
	Long: `
Backups of your Home Assistant system, which you can create,
restore, and delete using this command.`,
	Example: `
  ha backups
  ha backups new`,
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
	for idx, arg := range os.Args {
		if idx != 0 && (arg == "snapshot" || arg == "snap" || arg == "shot" || arg == "sn") {
			cmd.PrintErrf("The use of '%s' is deprecated, please use 'backups' instead!\n", arg)
		}
	}
	rootCmd.PersistentPreRun(cmd, args)
},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("backups")

		section := "backups"
		command := ""
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
	log.Debug("Init backups")
	// add cmd to root command
	rootCmd.AddCommand(backupsCmd)
}

package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var coreRebuildCmd = &cobra.Command{
	Use:     "rebuild",
	Aliases: []string{"rb", "reinstall"},
	Short:   "Rebuild the Home Assistant Core instance",
	Long: `
This command allows you to trigger a rebuild for your Home Assistant Core
instance running on your Home Assistant system.
Don't worry, this does not delete your config.`,
	Example: `
  ha core rebuild`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core rebuild")

		section := "core"
		command := "rebuild"
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
	coreCmd.AddCommand(coreRebuildCmd)
}

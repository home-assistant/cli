package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		ProgressSpinner.Start()
		resp, err := helper.GenericJSONPostTimeout(section, command, nil, helper.ContainerOperationTimeout)
		ProgressSpinner.Stop()
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	coreCmd.AddCommand(coreRebuildCmd)
}

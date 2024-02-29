package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osDataDiskWipeCmd = &cobra.Command{
	Use:     "wipe",
	Aliases: []string{"wipe", "reset", "erase"},
	Short:   "Wipe the Home Assistant Operating-System data partition",
	Long: `
This command will wipe all config for addons and Home Assistant and any locally
stored data in config, backups, media, etc. The machine will reboot during this.

After the reboot completes the latest stable version of Home Assistant and Supervisor
will be downloaded. Once the process is complete you will see onboarding, like
during initial setup.

The Operating System will not change during this. This includes OS level settings
such as network settings.

Please note, this command is limited due to security reasons, and will
only work on some locations. For example, the Operating System CLI.
`,
	Example: `
  ha os datadisk wipe
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os datadisk wipe")

		section := "os"
		command := "datadisk/wipe"

		if helper.AskForConfirmation(`
This will completely wipe the datadisk. This process is irreversible.
Are you sure you want to proceed?`, 0) {
			resp, err := helper.GenericJSONPost(section, command, nil)
			if err != nil {
				fmt.Println(err)
				ExitWithError = true
			} else {
				ExitWithError = !helper.ShowJSONResponse(resp)
			}
		}
	},
}

func init() {
	osDataDiskCmd.AddCommand(osDataDiskWipeCmd)
}

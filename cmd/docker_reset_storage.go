package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var dockerResetStorageCmd = &cobra.Command{
	Use:   "reset-storage",
	Short: "Reset the Docker storage",
	Long: `
This command wipes the complete Docker storage on the next reboot. All container
images, including Supervisor, Home Assistant Core, plugins and apps, are deleted
and downloaded again once the system is back up.

Configuration and data of Home Assistant and apps are kept, they are stored
outside of the Docker storage.

Internet connectivity is required for re-download of all the container images.
A reboot is required to apply the reset.
`,
	Example: `
  ha docker reset-storage
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("docker reset-storage", "args", args)

		section := "docker"
		command := "reset-storage"

		confirmed, err := helper.AskForConfirmation(`
This will delete all container images on the next reboot. Supervisor, Home
Assistant Core, plugins and apps are downloaded again afterwards, which
requires internet connectivity. Home Assistant and app data are kept.

Once confirmed, the reset will be applied on the next system reboot.

Are you sure you want to proceed?`, 0)

		if err != nil {
			cmd.PrintErrln("Aborted:", err)
			ExitWithError = true
			return
		}

		if confirmed {
			resp, err := helper.GenericJSONPost(section, command, nil)
			if err != nil {
				helper.PrintError(err)
				ExitWithError = true
			} else {
				ExitWithError = !helper.ShowJSONResponse(resp)
			}
		} else {
			cmd.PrintErrln("Aborted.")
			ExitWithError = true
		}
	},
}

func init() {
	dockerCmd.AddCommand(dockerResetStorageCmd)
}

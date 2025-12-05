package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dockerMigrateStorageCmd = &cobra.Command{
	Use:   "migrate-storage-driver [driver]",
	Short: "Migrate Docker storage driver",
	Long: `
This command schedules a Docker storage driver migration to be performed on the
next reboot. By default, it migrates to the (Containerd snapshotter) overlayfs.

The migration will be applied during the next system reboot. A reboot is required
to complete the migration.
`,
	Example: `
  ha docker migrate-storage-driver
  ha docker migrate-storage-driver overlayfs
`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return []string{"overlayfs"}, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("docker migrate-storage-driver")

		section := "docker"
		command := "migrate-storage-driver"

		// Default to overlayfs if no argument provided
		storageDriver := "overlayfs"
		if len(args) > 0 {
			storageDriver = args[0]
		}

		confirmed, err := helper.AskForConfirmation(`
This will schedule a Docker storage driver migration to "`+storageDriver+`".
Make sure to create a full Home Assistant backup before proceeding.

Internet connectivity is required for re-download of all the container images
and it is recommended to have at least 50% of free storage.

Once confirmed, the migration will be applied on the next system reboot.
Are you sure you want to proceed?`, 0)

		if err != nil {
			cmd.PrintErrln("Aborted:", err)
			ExitWithError = true
			return
		}

		if confirmed {
			options := map[string]any{
				"storage_driver": storageDriver,
			}
			resp, err := helper.GenericJSONPost(section, command, options)
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
	dockerCmd.AddCommand(dockerMigrateStorageCmd)
}

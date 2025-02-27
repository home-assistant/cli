package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var mountsOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Set options for mount manager",
	Long: `
Change value for options of mount manager in Supervisor.
`,
	Example: `
  ha mounts options --default-backup-mount my_share
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("mounts options")

		section := "mounts"
		command := "options"

		options := make(map[string]any)

		backupMount, err := cmd.Flags().GetString("default-backup-mount")
		if err == nil && cmd.Flags().Changed("default-backup-mount") {
			if backupMount == "" {
				options["default_backup_mount"] = nil
			} else {
				options["default_backup_mount"] = backupMount
			}
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	mountsOptionsCmd.Flags().String("default-backup-mount", "", "Set default backup mount")
	mountsOptionsCmd.Flags().Lookup("default-backup-mount").NoOptDefVal = ""
	mountsOptionsCmd.RegisterFlagCompletionFunc("default-backup-mount", backupsLocationsCompletions)

	mountsOptionsCmd.Flags().SetNormalizeFunc(func(set *pflag.FlagSet, name string) pflag.NormalizedName { // backwards compatibility
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	})

	mountsCmd.AddCommand(mountsOptionsCmd)
}

package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var backupOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allow to set options on backup manager",
	Long: `
This command allows you to set configuration options for backup manager.`,
	Example: `
  ha backups options --days-until-stale 60`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("backups options")

		section := "backups"
		command := "options"

		options := make(map[string]any)

		daysUntilStale, err := cmd.Flags().GetInt("days-until-stale")
		if daysUntilStale != 0 && err == nil && cmd.Flags().Changed("days-until-stale") {
			options["days_until_stale"] = daysUntilStale
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
	backupOptionsCmd.Flags().Int("days-until-stale", 30, "Days until backup considered stale")
	backupOptionsCmd.Flags().SetNormalizeFunc(func(set *pflag.FlagSet, name string) pflag.NormalizedName { // backwards compatibility
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	})
	backupOptionsCmd.RegisterFlagCompletionFunc("days-until-stale", cobra.NoFileCompletions)
	backupsCmd.AddCommand(backupOptionsCmd)
}

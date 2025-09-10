package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostDisksUsageCmd = &cobra.Command{
	Use:     "usage",
	Aliases: []string{"us", "use"},
	Short:   "Get default disk usage information",
	Long: `
This command provides information about the default disk usage on the host system
that Home Assistant is running on.`,
	Example: `
  ha host disks usage`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host disks usage")

		section := "host"
		command := "disks/default/usage"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	hostDisksCmd.AddCommand(hostDisksUsageCmd)
}

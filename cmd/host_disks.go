package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var hostDisksCmd = &cobra.Command{
	Use:     "disks",
	Aliases: []string{"disk"},
	Short:   "Get information about host disks and their usage",
	Long: `
The disks command lists the disks of the host system that Home Assistant is
running on, including the partitions which can be used for disk mounts. It also
provides access to other disk-related operations.`,
	Example: `
  ha host disks
  ha host disks usage`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("host disks", "args", args)

		section := "host"
		command := "disks"

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
	hostCmd.AddCommand(hostDisksCmd)
}

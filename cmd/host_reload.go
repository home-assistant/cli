package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"update", "refresh", "re"},
	Short:   "Reload information from the host machine",
	Long: `
This commands reload the information Home Assistant has on the hostmachine.
If some setting are changed outside of Home Assistant, this commands updates
the internals of Home Assistant.`,
	Example: `
  ha host reload`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host reload")

		section := "host"
		command := "reload"

		resp, err := helper.GenericJSONPost(section, command, nil)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	hostCmd.AddCommand(hostReloadCmd)
}

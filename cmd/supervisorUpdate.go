package cmd

import (
	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// supervisorUpdateCmd represents the update subcommand for supervisor
var supervisorUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		basePath := "supervisor"
		endpoint := "update"
		get := false
		serverOverride := ""

		log.WithFields(log.Fields{
			"basepath":       basePath,
			"endpoint":       endpoint,
			"serverOverride": serverOverride,
			"get":            get,
			"options":        supervisorOpts,
			"rawJSON":        supervisorRawJSON,
			"filter":         supervisorFilter,
		}).Debug("[CmdSupervisor]")
		helpers.ExecCommand(basePath, endpoint, serverOverride, get, supervisorOpts, supervisorFilter, supervisorRawJSON)
	},
}

func init() {
	supervisorCmd.AddCommand(supervisorUpdateCmd)
}

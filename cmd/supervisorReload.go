package cmd

import (
	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// supervisorReloadCmd represents the reload subcommand for supervisor
var supervisorReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		basePath := "supervisor"
		endpoint := "info"
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
	supervisorCmd.AddCommand(supervisorReloadCmd)
}

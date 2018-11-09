package cmd

import (
	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// supervisorLogsCmd represents the logs subcommand for supervisor
var supervisorLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		basePath := "supervisor"
		endpoint := "logs"
		get := true
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
	supervisorCmd.AddCommand(supervisorLogsCmd)
}

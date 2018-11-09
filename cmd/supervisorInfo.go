package cmd

import (
	"github.com/home-assistant/hassio-cli/command/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// supervisorInfoCmd represents the info subcommand for supervisor
var supervisorInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		basePath := "supervisor"
		endpoint := "info"
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
	supervisorCmd.AddCommand(supervisorInfoCmd)
}

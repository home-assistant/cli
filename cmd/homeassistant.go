package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// homeassistantCmd represents the homeassistant command when called without any subcommands
var homeassistantCmd = &cobra.Command{
	Use:     "homeassistant",
	Aliases: []string{"ha"},
	Short:   "homeassistant ",
}

func init() {
	log.Debug("Init homeassistant")
	// add subcommands
	homeassistantCmd.AddCommand(homeassistantInfoCmd)
	homeassistantCmd.AddCommand(homeassistantLogsCmd)
	homeassistantCmd.AddCommand(homeassistantCheckCmd)
	homeassistantCmd.AddCommand(homeassistantRestartCmd)
	homeassistantCmd.AddCommand(homeassistantStartCmd)
	homeassistantCmd.AddCommand(homeassistantStopCmd)

	// add cmd to root command
	rootCmd.AddCommand(homeassistantCmd)
}

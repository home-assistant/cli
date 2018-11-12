package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// homeassistantCmd represents the homeassistant command when called without any subcommands
var homeassistantCmd = &cobra.Command{
	Use:     "homeassistant",
	Aliases: []string{"ha"},
}

func init() {
	log.Debug("Init homeassistant")
	// add subcommands
	homeassistantCmd.AddCommand(homeassistantInfoCmd)

	// add cmd to root command
	rootCmd.AddCommand(homeassistantCmd)
}

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var homeassistantInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant info")
	},
}

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hardwareCmd = &cobra.Command{
	Use:     "hardware",
	Aliases: []string{"hw"},
}

func init() {
	log.Debug("Init hardware")

	// add cmd to root command
	rootCmd.AddCommand(hardwareCmd)
}

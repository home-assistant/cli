package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// hardwareCmd represents the hardware command
var hardwareCmd = &cobra.Command{
	Use:     "hardware",
	Aliases: []string{"hw"},
}

func init() {
	log.Debug("Init hardware")
	// add subcommands
	// TODO: add subcommand

	// add cmd to root command
	rootCmd.AddCommand(hardwareCmd)
}

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// addonsCmd represents the addons command
var addonsCmd = &cobra.Command{
	Use:     "addons",
	Aliases: []string{"ad"},
}

func init() {
	log.Debug("Init addons")
	// add subcommands
	// TODO: add subcommand

	// add cmd to root command
	rootCmd.AddCommand(addonsCmd)
}

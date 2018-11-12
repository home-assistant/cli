package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// supervisorCmd represents the supervisor command
var supervisorCmd = &cobra.Command{
	Use:     "supervisor",
	Aliases: []string{"su"},
}

func init() {
	log.Debug("Init supervisor")
	// add subcommands
	// TODO: add subcommand

	// add cmd to root command
	rootCmd.AddCommand(supervisorCmd)
}

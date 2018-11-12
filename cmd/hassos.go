package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// hassosCmd represents the hassos command
var hassosCmd = &cobra.Command{
	Use:     "hassos",
	Aliases: []string{"os"},
}

func init() {
	log.Debug("Init hassos")
	// add subcommands
	// TODO: add subcommand

	// add cmd to root command
	rootCmd.AddCommand(hassosCmd)
}

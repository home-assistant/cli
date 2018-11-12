package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"ho"},
}

func init() {
	log.Debug("Init host")
	// add subcommands
	// TODO: add subcommand

	// add cmd to root command
	rootCmd.AddCommand(hostCmd)
}

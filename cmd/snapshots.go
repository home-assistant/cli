package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// snapshotsCmd represents the snapshots command
var snapshotsCmd = &cobra.Command{
	Use:     "snapshots",
	Aliases: []string{"sa"},
}

func init() {
	log.Debug("Init snapshots")
	// add subcommands
	// TODO: add subcommand

	// add cmd to root command
	rootCmd.AddCommand(snapshotsCmd)
}

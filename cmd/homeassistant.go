package cmd

import (
	"github.com/spf13/cobra"
)

// homeassistantCmd represents the homeassistant command when called without any subcommands
var homeassistantCmd = &cobra.Command{
	Use:     "homeassistant",
	Aliases: []string{"ha"},
	Short:   "homeassistant ",
}

func init() {
	// add cmd to root command
	rootCmd.AddCommand(homeassistantCmd)
}

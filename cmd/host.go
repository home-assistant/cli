package cmd

import (
	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"ho"},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}

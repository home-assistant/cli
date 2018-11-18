package cmd

import (
	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"ho"},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}

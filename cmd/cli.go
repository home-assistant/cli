package cmd

import (
	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Get information, update or configure the Home Assistant cli backend",
	Long: `
The dns command allows you to manage the internal Home Assistant cli backend by
exposing commands to view, monitor, configure and control it.`,
	Example: `
  ha cli info
  ha cli update`,
}

func init() {
	rootCmd.AddCommand(cliCmd)
}

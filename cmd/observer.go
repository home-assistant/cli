package cmd

import (
	"github.com/spf13/cobra"
)

var observerCmd = &cobra.Command{
	Use:   "observer",
	Short: "Get information, update or configure the Home Assistant observer",
	Long: `
The observer command allows you to manage the internal Home Assistant observer by
exposing commands to view, monitor, configure and control it.`,
	Example: `
  ha observer info
  ha observer update`,
}

func init() {
	rootCmd.AddCommand(observerCmd)
}

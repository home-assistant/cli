package cmd

import (
	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:     "time",
	Aliases: []string{"ntp", "timedate"},
	Short:   "Get information or configure Home Assistant OS time settings",
	Long: `
The time command allows you to view and configure the Home Assistant OS network
time synchronization servers.`,
	Example: `
  ha time info
  ha time options --servers pool.ntp.org --fallback-servers time.google.com`,
}

func init() {
	rootCmd.AddCommand(timeCmd)
}

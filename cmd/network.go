package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"net"},
	Short:   "Network specific for updating, info and configuration imports",
	Long: `
The network command provides command line tools to control the host network that
Home Assistant is running on. It allows you to do things like change the
system network IP address, set connection options or join a Wi-Fi network.`,
	Example: `
  ha network info
  ha network interface options`,
}

func init() {
	log.Debug("Init network")
	rootCmd.AddCommand(networkCmd)
}

func ipMethodCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"static", "auto", "disabled"}, cobra.ShellCompDirectiveNoFileComp
}

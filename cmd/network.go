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
The network command provides commandline tools to control the host network that
Home Assistant is running on. It allows you do thing like change the
system network IP address or set connection options or join into a wifi.`,
	Example: `
  ha network info
  ha network interface options`,
}

func init() {
	log.Debug("Init network")
	rootCmd.AddCommand(networkCmd)
}

package cmd

import (
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:     "dns",
	Short:   "Get information, update or configure the Hass.io DNS server",
	Long: `
The dns command allows you to manage the internal Hass.io DNS server by
exposing commands to view, monitor, configure and control it.`,
	Example: `
  hassio dns logs
  hassio dns info
  hassio dns update`,
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}

package cmd

import (
	"github.com/spf13/cobra"
)

var securityCmd = &cobra.Command{
	Use:     "security",
	Aliases: []string{"secure", "sec"},
	Short:   "Get information and manage security functionality",
	Long: `
The security command allows you to manage the internal Home Assistant Security backend and
exposing commands to view, configure and control it.`,
	Example: `
  ha security info
  ha security options`,
}

func init() {
	rootCmd.AddCommand(securityCmd)
}

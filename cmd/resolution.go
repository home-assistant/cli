package cmd

import (
	"github.com/spf13/cobra"
)

var resolutionCmd = &cobra.Command{
	Use:     "resolution",
	Aliases: []string{"resolutions", "res"},
	Short:   "Resolution center of Supervisor, show issues and suggest solutions",
	Long: `
The Resolution center provide information about detected issue on the system. It allow also 
to give suggestion which can fix the issue and can executed by the suggestion command. Othewise
it's possible to dismiss issues or suggestion. It show also why a system show as not supported.`,
	Example: `
  ha resolution info
  ha resolution suggestion apply [ID]`,
}

func init() {
	// add cmd to root command
	rootCmd.AddCommand(resolutionCmd)
}

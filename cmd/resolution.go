package cmd

import (
	"github.com/spf13/cobra"
)

var resolutionCmd = &cobra.Command{
	Use:     "resolution",
	Aliases: []string{"resolutions", "res"},
	Short:   "Resolution center of Supervisor, show issues and suggest solutions",
	Long: `
The Resolution center provides information about detected issues on the system.
It also gives suggestions that can fix the issue and can be executed by the suggestion command.
It is possible to dismiss issues or suggestions. 
If you are running an unsupported system, the reasons for it will also show here`,
	Example: `
  ha resolution info
  ha resolution suggestion apply [ID]`,
}

func init() {
	// add cmd to root command
	rootCmd.AddCommand(resolutionCmd)
}

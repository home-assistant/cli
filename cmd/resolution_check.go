package cmd

import (
	"github.com/spf13/cobra"
)

var resolutionCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"checks", "test", "che", "ch"},
	Short:   "Check management by Resolution center",
	Long: `
This command allows to manage checks they runs by the system.`,
	Example: `
  ha resolution check options [slug]`,
}

func init() {
	resolutionCmd.AddCommand(resolutionCheckCmd)
}

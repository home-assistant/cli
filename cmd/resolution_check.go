package cmd

import (
	"github.com/spf13/cobra"
)

var resolutionCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"checks", "test", "che", "ch"},
	Short:   "Check management by the Resolution center",
	Long: `
This command allows to manage checks that are run by the system.`,
	Example: `
  ha resolution check options [slug]`,
}

func init() {
	resolutionCmd.AddCommand(resolutionCheckCmd)
}

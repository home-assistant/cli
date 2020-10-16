package cmd

import (
	"github.com/spf13/cobra"
)

var resolutionIssueCmd = &cobra.Command{
	Use:     "issue",
	Aliases: []string{"is", "trouble"},
	Short:   "Issues management reported by Resolution center",
	Long: `
This command allows dismissing issues reported by the system.`,
	Example: `
  ha resolution issue dismiss [id]`,
}

func init() {
	resolutionCmd.AddCommand(resolutionIssueCmd)
}

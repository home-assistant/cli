package cmd

import (
	"github.com/spf13/cobra"
)

var resolutionSuggestionCmd = &cobra.Command{
	Use:     "suggestion",
	Aliases: []string{"su", "solution"},
	Short:   "Suggestion management reported by Resolution center",
	Long: `
This command allow to dismiss or apply suggestion reported by the System.`,
	Example: `
  ha resolution suggestion dismiss [id]
  ha resolution suggestion apply [id]`,
}

func init() {
	resolutionCmd.AddCommand(resolutionSuggestionCmd)
}

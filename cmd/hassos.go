package cmd

import (
	"github.com/spf13/cobra"
)

// hassosCmd represents the hassos command
var hassosCmd = &cobra.Command{
	Use:     "hassos",
	Aliases: []string{"os"},
}

func init() {
	rootCmd.AddCommand(hassosCmd)
}

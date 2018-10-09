package hardware

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hardware",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hardware sub-command")
	},
}

// Execute represents the entrypoint for subcommands associated with "hardware"
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

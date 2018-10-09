package snapshots

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hassio",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snapshots subcommand")
	},
}

// Execute represents the entrypoint for snapshots subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

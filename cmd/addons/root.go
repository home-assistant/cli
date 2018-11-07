package addons

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "addons",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addons sub-command")
	},
}

// Execute represents the entrypoint for subcommands associated with "addons"
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

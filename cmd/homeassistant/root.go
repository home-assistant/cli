package homeassistant

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "homeassistant",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("homeassistant subcommand")
	},
}

// Execute represents the entrypoint for homeassistant subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package supervisor

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hassio",
	Short: "A brief description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("supervisor subcommand")
	},
}

// Execute represents the entrypoint to the supervisor subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

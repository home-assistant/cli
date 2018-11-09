package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// supervisorCmd represents the supervisor command
var supervisorCmd = &cobra.Command{
	Use:   "supervisor",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "No valid action detected.\n")
		os.Exit(3)
	},
}

func init() {
	rootCmd.AddCommand(supervisorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// supervisorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// supervisorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

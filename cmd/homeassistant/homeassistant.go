package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rawjson bool
var options string
var filter string

// homeassistantCmd represents the homeassistant command
var homeassistantCmd = &cobra.Command{
	Use:   "homeassistant",
	Short: "A brief description of your command",
	Aliases: []string{"ha"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("homeassistant called")
	},
}

func init() {
	rootCmd.AddCommand(homeassistantCmd)

	homeassistantCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// homeassistantCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package cmd

import (
	"fmt"
        "os"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
	"github.com/home-assistant/hassio-cli/command/helpers"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		basePath := "supervisor"
		endpoint := "info"
		get := true
		
	},
}

func init() {
	supervisor.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

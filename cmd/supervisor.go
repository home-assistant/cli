package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var supervisorOpts string
var supervisorRawJSON = false
var supervisorFilter string

// supervisorCmd represents the supervisor command
var supervisorCmd = &cobra.Command{
	Use:     "supervisor",
	Aliases: []string{"su"},
	Short:   "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "No valid action detected.\n")
		os.Exit(3)
	},
}

func init() {
	rootCmd.AddCommand(supervisorCmd)

	supervisorCmd.PersistentFlags().StringVarP(&supervisorOpts, "options", "o", supervisorOpts, "holds data for POST in format `key=val,key2=val2`")
	supervisorCmd.PersistentFlags().StringVarP(&supervisorFilter, "filter", "f", supervisorFilter, "properties to extract from returned data `prop1,prop2`")
	supervisorCmd.PersistentFlags().BoolVarP(&supervisorRawJSON, "rawjson", "j", supervisorRawJSON, "Returns the output in JSON format")
}

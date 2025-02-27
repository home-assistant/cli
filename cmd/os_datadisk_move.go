package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osDataDiskMoveCmd = &cobra.Command{
	Use:     "move [disk]",
	Aliases: []string{"migrate", "mov"},
	Short:   "Migrate Home Assistant Operating-System data partition",
	Long: `
This commands triggers an migration of the Home Assistant Operating-System
data partition to a new harddisk. The system reboots afterwards!
`,
	Example: `
  ha os datadisk move /dev/sda
`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if toComplete == "" {
			return []string{"/dev/"}, cobra.ShellCompDirectiveNoSpace
		}
		return nil, cobra.ShellCompDirectiveDefault
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os datadisk move")

		section := "os"
		command := "datadisk/move"
		options := make(map[string]any)

		options["device"] = args[0]

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	osDataDiskCmd.AddCommand(osDataDiskMoveCmd)
}

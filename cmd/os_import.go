package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var osImportCmd = &cobra.Command{
	Use:     "import",
	Aliases: []string{"im", "sync", "load"},
	Short:   "Import configurations from a USB stick",
	Long: `
This commands triggers an import action from a connected USB stick with
configuration to load for the Home Assistant Operating System.
`,
	Example: `
  ha os import
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("os import")

		section := "os"
		command := "config/sync"

		resp, err := helper.GenericJSONPost(section, command, nil)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	osCmd.AddCommand(osImportCmd)
}

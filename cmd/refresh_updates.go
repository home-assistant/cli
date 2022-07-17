package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var refreshUpdatesCmd = &cobra.Command{
	Use:     "refresh-updates",
	Aliases: []string{"refresh", "refresh_updates"},
	Short:   "Reload stores and version information",
	Long: `
This command reloads information about add-on repositories and fetches new version files.
	`,
	Example: `
  ha refresh-update
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("refresh_updates")

		section := "refresh_updates"
		command := ""

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
	rootCmd.AddCommand(refreshUpdatesCmd)
}

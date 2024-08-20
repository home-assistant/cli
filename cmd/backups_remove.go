package cmd

import (
	"fmt"

	"github.com/home-assistant/cli/client"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var backupsRemoveCmd = &cobra.Command{
	Use:     "remove [slug]",
	Aliases: []string{"delete", "del", "rem", "rm"},
	Short:   "Deletes a backup from disk",
	Long: `
Backups can take quite a bit of diskspace, this command allows you to
clean backups from disk.`,
	Example: `
  ha backups remove c1a07617`,
	ValidArgsFunction: backupsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("backups remove")

		section := "backups"
		command := "{slug}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.Delete(url)
		resp, err = client.GenericJSONErrorHandling(resp, err)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {

	backupsCmd.AddCommand(backupsRemoveCmd)
}

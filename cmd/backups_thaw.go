package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var backupsThawCmd = &cobra.Command{
	Use:     "thaw",
	Aliases: []string{"th"},
	Short:   "Thaw supervisor after an external backup",
	Long: `
End a freeze initiated by the freeze command after an external backup or snapshot
has completed.`,
	Example: `
  ha backups thaw`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("backups thaw")

		section := "backups"
		command := "thaw"

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
	backupsCmd.AddCommand(backupsThawCmd)
}

package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allow to set options on host system",
	Long: `
This command allows you to set configuration options on the host system that 
your Home Assistant is running on.`,
	Example: `
  ha host options --hostname homeassistant.local`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host options")

		section := "host"
		command := "options"

		var options map[string]any

		hostname, _ := cmd.Flags().GetString("hostname")
		if hostname != "" {
			options = map[string]any{"hostname": hostname}
		}

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
	hostOptionsCmd.Flags().StringP("hostname", "", "", "Hostname to set")
	hostOptionsCmd.RegisterFlagCompletionFunc("hostname", cobra.NoFileCompletions)
	hostCmd.AddCommand(hostOptionsCmd)
}

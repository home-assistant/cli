package cmd

import (
	"fmt"
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostLogsCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"log", "lg"},
	Short:   "View the log output of the host systemd journal logs",
	Long: `
Allows you to look at the systemd journal on the host to see logs
across services and boots.
`,
	Example: `
  ha host logs
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host logs")

		section := "host"

		request, err := processLogsFlags(section, cmd)

		if err != nil {
			fmt.Printf("Error: %v", err)
			ExitWithError = true
			return
		}

		identifier, _ := cmd.Flags().GetString("identifier")
		if len(identifier) > 0 {
			if strings.HasSuffix(request.URL, "/follow") {
				// We can safely do this because "/follow" will be always at the end
				request.URL = strings.Replace(request.URL, "/follow", "/identifiers/{identifier}/follow", 1)
			} else {
				request.URL += "/identifiers/{identifier}"
			}
		}
		request.SetPathParam("identifier", identifier)

		resp, err := request.Send()

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		ExitWithError = !helper.StreamTextResponse(resp)
	},
}

func init() {
	addLogsFlags(hostLogsCmd)

	hostLogsCmd.Flags().StringP("identifier", "t", "", "Show entries with the specified syslog identifier")

	hostCmd.AddCommand(hostLogsCmd)
}

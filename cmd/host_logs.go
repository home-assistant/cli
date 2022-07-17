package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostLogsCmd = &cobra.Command{
	Use:     "logs",
	Aliases: []string{"log", "lg"},
	Short:   "View the log output of the host kernel logs/dmesg",
	Long: `
Allowing you to look at the log output generated by the host kernel.
`,
	Example: `
  ha host logs
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host logs")

		section := "host"
		command := "logs"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetRequest()
		resp, err := request.SetHeader("Accept", "text/plain").Get(url)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			fmt.Println(resp.String())
		}
	},
}

func init() {
	hostCmd.AddCommand(hostLogsCmd)
}

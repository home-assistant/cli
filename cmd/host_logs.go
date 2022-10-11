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
		command := "logs"

		listIdentifiers, _ := cmd.Flags().GetBool("list-identifiers")
		listBoots, _ := cmd.Flags().GetBool("list-boots")

		if listIdentifiers || listBoots {
			if listIdentifiers {
				command += "/identifiers"
			} else {
				command += "/boots"
			}

			resp, err := helper.GenericJSONGet(section, command)
			if err != nil {
				fmt.Println(err)
				ExitWithError = true
			} else {
				ExitWithError = !helper.ShowJSONResponse(resp)
			}

			return
		}

		identifier, _ := cmd.Flags().GetString("identifier")
		boot, _ := cmd.Flags().GetString("boot")
		if len(boot) > 0 {
			command += "/boots/{boot}"
		}
		if len(identifier) > 0 {
			command += "/identifiers/{identifier}"
		}

		follow, _ := cmd.Flags().GetBool("follow")
		if follow {
			command += "/follow"
		}

		url, err := helper.URLHelper(section, command)

		if err != nil {
			fmt.Printf("Error: %v", err)
			ExitWithError = true
			return
		}

		/* Disable timeouts to allow following forever */
		request := helper.GetRequestTimeout(0).SetHeader("Accept", "text/plain").SetDoNotParseResponse(true)

		lines, _ := cmd.Flags().GetInt32("lines")
		if lines > 0 {
			request.SetHeader("Range", fmt.Sprintf("entries=:%d:", -(lines-1)))
		}

		if err != nil {
			fmt.Printf("Error: %v", err)
			ExitWithError = true
			return
		}

		request.SetPathParam("identifier", identifier)
		request.SetPathParam("boot", boot)

		resp, err := request.Get(url)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		ExitWithError = !helper.StreamTextResponse(resp)
	},
}

func init() {
	hostLogsCmd.Flags().BoolP("follow", "f", false, "Continuously print new log entries")
	hostLogsCmd.Flags().Int32P("lines", "n", 0, "Number of log entries to show")
	hostLogsCmd.Flags().StringP("identifier", "t", "", "Show entries with the specified syslog identifier")
	hostLogsCmd.Flags().StringP("boot", "b", "", "Logs of particular boot ID")
	hostLogsCmd.Flags().Bool("list-identifiers", false, "Show all syslog identifiers")
	hostLogsCmd.Flags().Bool("list-boots", false, "Show all boot IDs by offset")

	hostLogsCmd.Flags().Lookup("follow").NoOptDefVal = "true"
	hostLogsCmd.Flags().Lookup("list-identifiers").NoOptDefVal = "true"
	hostLogsCmd.Flags().Lookup("list-boots").NoOptDefVal = "true"

	hostCmd.AddCommand(hostLogsCmd)
}

package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkScanCmd = &cobra.Command{
	Use:     "scan [interface]",
	Aliases: []string{"accesspoints", "wifi"},
	Short:   "Scan for Access Points on a wireless interface.",
	Long: `
Scan for Access Points on a specific wireless interface.
This function works only on a wireless interface!
`,
	Example: `
  ha network scan wlan0
`,
	ValidArgsFunction: networkInterfaceCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("network scan")

		section := "network"
		command := "interface/{interface}/accesspoints"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		request := helper.GetJSONRequest()

		inet := args[0]

		request.SetPathParams(map[string]string{
			"interface": inet,
		})

		ProgressSpinner.Start()
		resp, err := request.Get(url)
		ProgressSpinner.Stop()

		// returns 200 OK or 400, everything else is wrong
		if err == nil {
			if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
				err = errors.New("Unexpected server response")
				log.Error(err)
			} else if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
				err = errors.New("API did not return a JSON response")
				log.Error(err)
			}
		}

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	networkCmd.AddCommand(networkScanCmd)
}

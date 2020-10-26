package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var networkUpdateCmd = &cobra.Command{
	Use:     "update [interface]",
	Aliases: []string{"set", "up"},
	Short:   "Update settings of a network interface",
	Long: `
Update network interface settings of a specific adapter.
`,
	Example: `
  ha network update eth0 --method dhcp
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("network update")

		section := "network"
		command := "interface/{interface}/update"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		options := make(map[string]interface{})

		request := helper.GetJSONRequest()

		inet := args[0]

		request.SetPathParams(map[string]string{
			"interface": inet,
		})

		for _, value := range []string{
			"address",
			"gateway",
			"method",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[value] = val
			}
		}

		dns, err := cmd.Flags().GetStringArray("dns")
		if len(dns) >= 1 && err == nil {
			options["dns"] = dns
		}

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		resp, err := request.Post(url)

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
	networkUpdateCmd.Flags().StringP("address", "a", "", "IP address for the interface in the 192.168.1.5/24")
	networkUpdateCmd.Flags().StringP("gateway", "g", "", "The gateway the interface should use")
	networkUpdateCmd.Flags().StringP("method", "m", "", "Method: static|dhcp")
	networkUpdateCmd.Flags().StringArrayP("dns", "d", []string{}, "Upstream DNS servers to use. Use multiple times for multiple servers.")
	networkCmd.AddCommand(networkUpdateCmd)
}

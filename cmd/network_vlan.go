package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkVlanCmd = &cobra.Command{
	Use:     "vlan [interface] [id]",
	Aliases: []string{},
	Short:   "Create a new VLAN on an ethernet interface",
	Long: `
Create a new VLAN on an ethernet interface. It allows setting an initial IP config.
This function works only on an ethernet interface!
`,
	Example: `
  ha network vlan eth0 10 --ipv4-method auto --ipv6-method disabled
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("network vlan")

		section := "network"
		command := "interface/{interface}/vlan/{vlan}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		options := make(map[string]interface{})

		request := helper.GetJSONRequest()

		inet := args[0]
		vlan := args[1]

		request.SetPathParams(map[string]string{
			"interface": inet,
			"vlan":      vlan,
		})

		// IP configs
		helperIpConfig("ipv4", cmd, options)
		helperIpConfig("ipv6", cmd, options)

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

	networkVlanCmd.Flags().StringArray("ipv4-address", []string{}, "IPv4 address for the interface in the 192.168.1.5/24")
	networkVlanCmd.Flags().String("ipv4-gateway", "", "The IPv4 gateway the interface should use")
	networkVlanCmd.Flags().String("ipv4-method", "", "Method on IPv4: static|auto|disabled")
	networkVlanCmd.Flags().StringArray("ipv4-nameserver", []string{}, "Upstream DNS servers to use for IPv4. Use multiple times for multiple servers.")

	networkVlanCmd.Flags().StringArray("ipv6-address", []string{}, "IPv6 address for the interface in the 2001:0db8:85a3:0000:0000:8a2e:0370:7334/64")
	networkVlanCmd.Flags().String("ipv6-gateway", "", "The IPv6 gateway the interface should use")
	networkVlanCmd.Flags().String("ipv6-method", "", "Method on IPv6: static|auto|disabled")
	networkVlanCmd.Flags().StringArray("ipv6-nameserver", []string{}, "Upstream DNS servers to use for IPv6. Use multiple times for multiple servers.")

	networkVlanCmd.RegisterFlagCompletionFunc("ipv4-address", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv4-gateway", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv4-method", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv4-nameserver", cobra.NoFileCompletions)

	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-address", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-gateway", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-method", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-nameserver", cobra.NoFileCompletions)

	networkCmd.AddCommand(networkVlanCmd)
}

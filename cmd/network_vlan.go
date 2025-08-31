package cmd

import (
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
	ValidArgsFunction: networkInterfaceCompletions,
	Args:              cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("network vlan")

		section := "network"
		command := "interface/{interface}/vlan/{vlan}"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		options := make(map[string]any)

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

		// mDNS / LLMNR
		helperMdnsConfig(cmd, options)

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		resp, err := request.Post(url)
		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {

	networkVlanCmd.Flags().StringArray("ipv4-address", []string{}, "IPv4 address for the interface in CIDR notation (e.g. 192.168.1.5/24)")
	networkVlanCmd.Flags().String("ipv4-gateway", "", "The IPv4 gateway the interface should use")
	networkVlanCmd.Flags().String("ipv4-method", "", "Method on IPv4: static|auto|disabled")
	networkVlanCmd.Flags().StringArray("ipv4-nameserver", []string{}, "IPv4 address of upstream DNS servers. Use multiple times for multiple servers.")

	networkVlanCmd.Flags().StringArray("ipv6-address", []string{}, "IPv6 address for the interface in CIDR notation (e.g. 2001:db8:85a3::8a2e:370:7334/64)")
	networkVlanCmd.Flags().String("ipv6-gateway", "", "The IPv6 gateway the interface should use")
	networkVlanCmd.Flags().String("ipv6-method", "", "Method on IPv6: static|auto|disabled")
	networkVlanCmd.Flags().String("ipv6-addr-gen-mode", "", "IPv6 address generation mode: eui64|stable-privacy|default-or-eui64|default")
	networkVlanCmd.Flags().String("ipv6-privacy", "", "IPv6 privacy extensions: disabled|enabled-prefer-public|enabled|default")
	networkVlanCmd.Flags().StringArray("ipv6-nameserver", []string{}, "IPv6 address for upstream DNS servers. Use multiple times for multiple servers.")

	networkVlanCmd.Flags().String("mdns", "", "mDNS mode: default|off|resolve|announce")
	networkVlanCmd.Flags().String("llmnr", "", "LLMNR mode: default|off|resolve|announce")

	networkVlanCmd.RegisterFlagCompletionFunc("ipv4-address", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv4-gateway", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv4-method", ipMethodCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv4-nameserver", cobra.NoFileCompletions)

	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-address", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-gateway", cobra.NoFileCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-method", ipMethodCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-addr-gen-mode", ipAddrGenModeCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-privacy", ip6PrivacyCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("ipv6-nameserver", cobra.NoFileCompletions)

	networkVlanCmd.RegisterFlagCompletionFunc("mdns", mdnsCompletions)
	networkVlanCmd.RegisterFlagCompletionFunc("llmnr", mdnsCompletions)

	networkCmd.AddCommand(networkVlanCmd)
}

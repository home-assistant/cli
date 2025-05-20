package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkUpdateCmd = &cobra.Command{
	Use:     "update [interface]",
	Aliases: []string{"set", "up"},
	Short:   "Update settings of a network interface",
	Long: `
Update network interface settings of a specific adapter.
`,
	Example: `
  ha network update eth0 --ipv4-method auto --ipv6-method disabled
`,
	ValidArgsFunction: networkInterfaceCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("network update")

		section := "network"
		command := "interface/{interface}/update"

		url, err := helper.URLHelper(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		options := make(map[string]any)

		request := helper.GetJSONRequest()

		inet := args[0]

		request.SetPathParams(map[string]string{
			"interface": inet,
		})

		// IP configs
		helperIpConfig("ipv4", cmd, options)
		helperIpConfig("ipv6", cmd, options)

		// Wifi
		helperWifiConfig(cmd, options)

		disabled, err := cmd.Flags().GetBool("disabled")
		if err == nil {
			options["enabled"] = !disabled
		}

		log.WithField("options", options).Debug("Request body")
		request.SetBody(options)

		resp, err := request.Post(url)
		resp, err = helper.GenericJSONErrorHandling(resp, err)

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {

	networkUpdateCmd.Flags().StringArray("ipv4-address", []string{}, "IPv4 address for the interface in the 192.168.1.5/24")
	networkUpdateCmd.Flags().String("ipv4-gateway", "", "The IPv4 gateway the interface should use")
	networkUpdateCmd.Flags().String("ipv4-method", "", "Method on IPv4: static|auto|disabled")
	networkUpdateCmd.Flags().StringArray("ipv4-nameserver", []string{}, "IPv4 address of upstream DNS servers. Use multiple times for multiple servers.")

	networkUpdateCmd.Flags().StringArray("ipv6-address", []string{}, "IPv6 address for the interface in the 2001:0db8:85a3:0000:0000:8a2e:0370:7334/64")
	networkUpdateCmd.Flags().String("ipv6-gateway", "", "The IPv6 gateway the interface should use")
	networkUpdateCmd.Flags().String("ipv6-method", "", "Method on IPv6: static|auto|disabled")
	networkUpdateCmd.Flags().String("ipv6-addr-gen-mode", "", "IPv6 address generation mode: eui64|stable-privacy")
	networkUpdateCmd.Flags().String("ipv6-privacy", "", "IPv6 privacy extensions: disabled|enabled-prefer-public|enabled")
	networkUpdateCmd.Flags().StringArray("ipv6-nameserver", []string{}, "IPv6 address for upstream DNS servers. Use multiple times for multiple servers.")

	networkUpdateCmd.Flags().String("wifi-mode", "", "Wifi mode: infrastructure, adhoc, mesh or ap")
	networkUpdateCmd.Flags().String("wifi-ssid", "", "SSID for wifi connection")
	networkUpdateCmd.Flags().String("wifi-auth", "", "Used authentication: open, wep, wpa-psk")
	networkUpdateCmd.Flags().String("wifi-psk", "", "Shared authentication key for wep or wpa")

	networkUpdateCmd.Flags().BoolP("disabled", "e", false, "Disable interface")

	networkUpdateCmd.RegisterFlagCompletionFunc("ipv4-address", cobra.NoFileCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("ipv4-gateway", cobra.NoFileCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("ipv4-method", ipMethodCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("ipv4-nameservers", cobra.NoFileCompletions)

	networkUpdateCmd.RegisterFlagCompletionFunc("ipv6-address", cobra.NoFileCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("ipv6-gateway", cobra.NoFileCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("ipv6-method", ipMethodCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("ipv6-addr-gen-mode", ipAddrGenModeCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("ipv6-privacy", ip6PrivacyCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("ipv6-nameservers", cobra.NoFileCompletions)

	networkUpdateCmd.RegisterFlagCompletionFunc("wifi-mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"infrastructure", "adhoc", "mesh", "ap"}, cobra.ShellCompDirectiveNoFileComp
	})
	networkUpdateCmd.RegisterFlagCompletionFunc("wifi-ssid", cobra.NoFileCompletions)
	networkUpdateCmd.RegisterFlagCompletionFunc("wifi-auth", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"open", "wep", "wpa-psk"}, cobra.ShellCompDirectiveNoFileComp
	})
	networkUpdateCmd.RegisterFlagCompletionFunc("wifi-psk", cobra.NoFileCompletions)

	networkUpdateCmd.RegisterFlagCompletionFunc("disabled", boolCompletions)

	networkCmd.AddCommand(networkUpdateCmd)
}

type NetworkArg struct {
	Arg     string
	ApiKey  string
	IsArray bool
}

func parseNetworkArgs(cmd *cobra.Command, args []NetworkArg) map[string]any {
	networkConfig := make(map[string]any)
	for _, arg := range args {
		var val any
		var err error
		var changed bool

		if arg.IsArray {
			val, err = cmd.Flags().GetStringArray(arg.Arg)
			changed = len(val.([]string)) > 0
		} else {
			val, err = cmd.Flags().GetString(arg.Arg)
			changed = val.(string) != ""
		}

		if err == nil && changed && cmd.Flags().Changed(arg.Arg) {
			networkConfig[arg.ApiKey] = val
		}
	}
	return networkConfig
}

func helperIpConfig(version string, cmd *cobra.Command, options map[string]any) {
	args := []NetworkArg{
		{Arg: version + "-gateway", ApiKey: "gateway"},
		{Arg: version + "-method", ApiKey: "method"},
		{Arg: version + "-addr-gen-mode", ApiKey: "addr_gen_mode"},
		{Arg: version + "-privacy", ApiKey: "ip6_privacy"},
		{Arg: version + "-address", ApiKey: "address", IsArray: true},
		{Arg: version + "-nameserver", ApiKey: "nameservers", IsArray: true},
	}

	ipConfig := parseNetworkArgs(cmd, args)
	if len(ipConfig) > 0 {
		options[version] = ipConfig
	}
}

func helperWifiConfig(cmd *cobra.Command, options map[string]any) {
	args := []NetworkArg{
		{Arg: "wifi-mode", ApiKey: "mode"},
		{Arg: "wifi-ssid", ApiKey: "ssid"},
		{Arg: "wifi-auth", ApiKey: "auth"},
		{Arg: "wifi-psk", ApiKey: "psk"},
	}

	wifiConfig := parseNetworkArgs(cmd, args)
	if len(wifiConfig) > 0 {
		options["wifi"] = wifiConfig
	}
}

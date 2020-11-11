package cmd

import (
	"errors"
	"fmt"
	"strings"

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
  ha network update eth0 --ipv4-method dhcp --ipv6-method disabled
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

		// IP configs
		helperIpConfig("ipv4", cmd, options)
		helperIpConfig("ipv6", cmd, options)

		// Wifi
		helperWifiConfig(cmd, options)

		enabled, err := cmd.Flags().GetBool("enabled")
		if err != nil {
			options["enabled"] = enabled
		}

		log.WithField("options", options).Debug("Request body")
		request.SetBody(options)

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

	networkUpdateCmd.Flags().StringArray("ip4-address", []string{}, "IPv4 address for the interface in the 192.168.1.5/24")
	networkUpdateCmd.Flags().String("ipv4-gateway", "", "The IPv4 gateway the interface should use")
	networkUpdateCmd.Flags().String("ipv4-method", "", "Method on IPv4: static|dhcp|disabled")
	networkUpdateCmd.Flags().StringArray("ipv4-nameserver", []string{}, "Upstream DNS servers to use for IPv4. Use multiple times for multiple servers.")

	networkUpdateCmd.Flags().StringArray("ip6-address", []string{}, "IPv6 address for the interface in the 2001:0db8:85a3:0000:0000:8a2e:0370:7334/64")
	networkUpdateCmd.Flags().String("ipv6-gateway", "", "The IPv6 gateway the interface should use")
	networkUpdateCmd.Flags().String("ipv6-method", "", "Method on IPv6: static|dhcp|disabled")
	networkUpdateCmd.Flags().StringArray("ipv6-nameserver", []string{}, "Upstream DNS servers to use for IPv6. Use multiple times for multiple servers.")

	networkUpdateCmd.Flags().String("wifi-mode", "", "Wifi mode: infrastructure, adhoc, mesh or ap")
	networkUpdateCmd.Flags().String("wifi-ssid", "", "SSID for wifi connection")
	networkUpdateCmd.Flags().String("wifi-auth", "", "Used authentication: open, wep, wpa-psk")
	networkUpdateCmd.Flags().String("wifi-psk", "", "Shared authentication key for wep or wpa")

	networkUpdateCmd.Flags().BoolP("enabled", "e", true, "Enable or disable interface")
	networkCmd.AddCommand(networkUpdateCmd)
}

func helperIpConfig(version string, cmd *cobra.Command, options map[string]interface{}) {
	ipConfig := make(map[string]interface{})

	for _, value := range []string{
		version + "-gateway",
		version + "-method",
	} {
		val, err := cmd.Flags().GetString(value)
		if val != "" && err == nil && cmd.Flags().Changed(value) {
			ipConfig[strings.Split("-", value)[1]] = val
		}
	}

	for _, value := range []string{
		version + "-address",
		version + "-nameservers",
	} {
		val, err := cmd.Flags().GetStringArray(value)
		if len(val) >= 1 && err == nil && cmd.Flags().Changed(value) {
			ipConfig[strings.Split("-", value)[1]] = val
		}
	}

	if len(ipConfig) > 0 {
		options[version] = ipConfig
	}
}

func helperWifiConfig(cmd *cobra.Command, options map[string]interface{}) {
	wifiConfig := make(map[string]interface{})

	for _, value := range []string{
		"wifi-mode",
		"wifi-ssid",
		"wifi-auth",
		"wifi-psk",
	} {
		val, err := cmd.Flags().GetString(value)
		if val != "" && err == nil && cmd.Flags().Changed(value) {
			wifiConfig[strings.Split("-", value)[1]] = val
		}
	}

	if len(wifiConfig) > 0 {
		options["wifi"] = wifiConfig
	}
}

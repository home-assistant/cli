package cmd

import (
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"net"},
	Short:   "Network specific for updating, info and configuration imports",
	Long: `
The network command provides command line tools to control the host network that
Home Assistant is running on. It allows you to do things like change the
system network IP address, set connection options or join a Wi-Fi network.`,
	Example: `
  ha network info
  ha network interface options`,
}

func init() {
	log.Debug("Init network")
	rootCmd.AddCommand(networkCmd)
}

func ipMethodCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"static", "auto", "disabled"}, cobra.ShellCompDirectiveNoFileComp
}

func ipAddrGenModeCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"eui64", "stable-privacy", "default-or-eui64", "default"}, cobra.ShellCompDirectiveNoFileComp
}

func ip6PrivacyCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"disabled", "enabled-prefer-public", "enabled", "default"}, cobra.ShellCompDirectiveNoFileComp
}

func mdnsCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"default", "off", "resolve", "announce"}, cobra.ShellCompDirectiveNoFileComp
}

func networkInterfaceCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("network", "info")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["interfaces"] != nil {
		if ifaces, ok := data.Data["interfaces"].([]any); ok {
			for _, iface := range ifaces {
				var m map[string]any
				if m, ok = iface.(map[string]any); !ok {
					continue
				}
				var s string
				switch cmd.Name() {
				case "scan":
					var b bool
					if b, ok = m["enabled"].(bool); !ok || !b {
						continue
					}
					if s, ok = m["type"].(string); !ok || s != "wireless" {
						continue
					}
				case "vlan":
					if s, ok = m["type"].(string); !ok || s != "ethernet" {
						continue
					}
				}
				if s, ok = m["interface"].(string); !ok || s == "" {
					continue
				}
				ret = append(ret, s)
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}

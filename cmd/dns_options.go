package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dnsOptionsCmd = &cobra.Command{
	Use:     "options",
	Aliases: []string{"option", "opt", "opts", "op"},
	Short:   "Allow to set options for the internal Home Assistant DNS server",
	Long: `
This command allows you to set configuration options for the internally
running Home Assistant DNS server.
`,
	Example: `
  ha dns options --servers dns://8.8.8.8 --servers dns://1.1.1.1
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("dns options")

		section := "dns"
		command := "options"

		options := make(map[string]interface{})

		servers, err := cmd.Flags().GetStringArray("servers")
		log.WithField("servers", servers).Debug("servers")

		if len(servers) >= 1 && err == nil {
			options["servers"] = servers
		}

		data, err := cmd.Flags().GetBool("fallback")
		if err == nil && cmd.Flags().Changed("fallback") {
			options["fallback"] = data
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dnsOptionsCmd.Flags().StringArrayP("servers", "r", []string{}, "Upstream DNS servers to use. Use multiple times for multiple servers.")
	dnsOptionsCmd.Flags().BoolP("fallback", "", true, "Enable/Disable fallback DNS (Cloudflare DoT)")

	dnsOptionsCmd.Flags().Lookup("fallback").NoOptDefVal = "true"

	dnsCmd.AddCommand(dnsOptionsCmd)
}

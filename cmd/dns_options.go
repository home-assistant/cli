package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		base := viper.GetString("endpoint")

		options := make(map[string]interface{})

		servers, err := cmd.Flags().GetStringArray("servers")
		log.WithField("servers", servers).Debug("servers")

		if len(servers) >= 1 && err == nil {
			options["servers"] = servers
		}

		resp, err := helper.GenericJSONPost(base, section, command, options)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}

		return
	},
}

func init() {
	dnsOptionsCmd.Flags().StringArrayP("servers", "r", []string{}, "Upstream DNS servers to use. Use multiple times for multiple servers.")
	dnsCmd.AddCommand(dnsOptionsCmd)
}

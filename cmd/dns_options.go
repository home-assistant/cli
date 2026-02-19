package cmd

import (
	"log/slog"

	helper "github.com/home-assistant/cli/client"
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
  ha dns options --search-domains example.com --search-domains corp.example.com
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("dns options", "args", args)

		section := "dns"
		command := "options"

		options := make(map[string]any)

		servers, err := cmd.Flags().GetStringArray("servers")
		slog.Debug("servers", "servers", servers)

		if len(servers) >= 1 && err == nil {
			options["servers"] = servers
		}

		searchDomains, err := cmd.Flags().GetStringArray("search-domains")
		log.WithField("search-domains", searchDomains).Debug("search-domains")

		if cmd.Flags().Changed("search-domains") && err == nil {
			options["search_domains"] = searchDomains
		}

		data, err := cmd.Flags().GetBool("fallback")
		if err == nil && cmd.Flags().Changed("fallback") {
			options["fallback"] = data
		}

		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dnsOptionsCmd.Flags().StringArrayP("servers", "r", []string{}, "Upstream DNS servers to use. Use multiple times for multiple servers.")
	dnsOptionsCmd.Flags().StringArray("search-domains", []string{}, "DNS search domains written to container resolv.conf. Use multiple times for multiple domains.")
	dnsOptionsCmd.Flags().BoolP("fallback", "", true, "Enable/Disable fallback DNS (Cloudflare DoT)")

	dnsOptionsCmd.Flags().Lookup("fallback").NoOptDefVal = "true"

	dnsOptionsCmd.RegisterFlagCompletionFunc("servers", cobra.NoFileCompletions)
	dnsOptionsCmd.RegisterFlagCompletionFunc("search-domains", cobra.NoFileCompletions)
	dnsOptionsCmd.RegisterFlagCompletionFunc("fallback", boolCompletions)

	dnsCmd.AddCommand(dnsOptionsCmd)
}

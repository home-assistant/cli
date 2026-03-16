package cmd

import (
	"fmt"
	"net/http"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var coreReloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"refresh"},
	Short:   "Reload the Home Assistant Core configuration",
	Long: `
Reload the Home Assistant Core YAML configuration without restarting.
Use --component to reload a specific component (e.g., automations, scripts).
Without --component, reloads the core configuration.`,
	Example: `
  ha core reload
  ha core reload --component automation
  ha core reload --component script
  ha core reload --component scene
  ha core reload --component group`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("core reload")

		component, _ := cmd.Flags().GetString("component")

		var section, command string
		if component != "" {
			section = "core/api/services/" + component
			command = "reload"
		} else {
			section = "core/api/services/homeassistant"
			command = "reload_all"
		}

		url, err := helper.URLHelper(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		resp, err := helper.GetRequest().Post(url)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		if resp.StatusCode() != http.StatusOK {
			if component != "" && resp.StatusCode() == http.StatusBadRequest {
				helper.PrintErrorString(fmt.Sprintf(
					"component %q not found or does not support reload", component))
			} else {
				helper.PrintErrorString(fmt.Sprintf("unexpected response (status: %d)", resp.StatusCode()))
			}
			ExitWithError = true
		} else {
			fmt.Println("Command completed successfully.")
		}
	},
}

func init() {
	coreReloadCmd.Flags().StringP("component", "c", "",
		"Specific component to reload (e.g., automation, script, scene, group)")
	coreReloadCmd.RegisterFlagCompletionFunc("component",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"automation",
				"script",
				"scene",
				"group",
				"input_boolean",
				"input_number",
				"input_select",
				"input_text",
				"input_datetime",
				"input_button",
				"timer",
				"counter",
				"template",
				"zone",
				"schedule",
				"person",
			}, cobra.ShellCompDirectiveNoFileComp
		})

	coreCmd.AddCommand(coreReloadCmd)
}

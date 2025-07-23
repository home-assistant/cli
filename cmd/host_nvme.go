package cmd

import (
	"strings"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostNvmeCmd = &cobra.Command{
	Use:   "nvme",
	Short: "Manage NVMe devices available on the host system",
	Long: `
Show information about or manage NVMe devices available on the host system.`,
	Example: `
  ha host nvme status /dev/nvme0n1`,
}

func init() {
	hostCmd.AddCommand(hostNvmeCmd)
}

func nvmeDeviceCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("host", "info")
	if err != nil || !resp.IsSuccess() {
		log.WithError(err).Debug("Failed to fetch host info for NVMe completion")
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["nvme_devices"] != nil {
		if devices, ok := data.Data["nvme_devices"].([]any); ok {
			for _, dev := range devices {
				var m map[string]any
				if m, ok = dev.(map[string]any); !ok {
					continue
				}
				var id string
				if id, ok = m["id"].(string); !ok || id == "" {
					continue
				}

				entry := id
				var path string
				if path, ok = m["path"].(string); ok && path != "" {
					entry += "\t" + path
				}
				if toComplete == "" || strings.HasPrefix(id, toComplete) {
					ret = append(ret, entry)
				}
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}

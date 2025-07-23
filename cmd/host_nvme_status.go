package cmd

import (
	"fmt"
	"net/url"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostNvmeStatusCmd = &cobra.Command{
	Use:   "status [device]",
	Short: "Show NVMe status for a device",
	Long: `
Show NVMe status for a specific device or for the current datadisk (if using
Home Assistant Operating System and the current datadisk is an NVMe device).
`,
	Example: `
  ha host nvme status
  ha host nvme status /dev/nvme0n1
`,
	ValidArgsFunction: nvmeDeviceCompletions,
	Args:              cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("host nvme status")

		section := "host"
		command := "nvme/status"
		if len(args) > 0 {
			device := url.PathEscape(args[0])
			command = fmt.Sprintf("nvme/%s/status", device)
		}

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	hostNvmeCmd.AddCommand(hostNvmeStatusCmd)
}

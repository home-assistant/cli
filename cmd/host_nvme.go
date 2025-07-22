package cmd

import (
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

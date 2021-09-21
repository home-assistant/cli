package cmd

import (
	"github.com/spf13/cobra"
)

var osDataDiskCmd = &cobra.Command{
	Use:     "datadisk",
	Aliases: []string{"data"},
	Short:   "Operating System DataDisk feature for managing data partition",
	Long: `
This command set is specifically designed for the Home Assistant Operating System
and only works on those systems. It provides an interface to get information
or migrating the data partition of the system to a different harddrive.`,
	Example: `
  ha os datadisk list
  ha os datadisk move /dev/sda1`,
}

func init() {
	osCmd.AddCommand(osDataDiskCmd)
}

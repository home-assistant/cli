package cmd

import (
	"github.com/spf13/cobra"
)

var factoryResetCmd = &cobra.Command{
	Use:   "factory-reset",
	Short: "Reset Home Assistant to factory settings (alias for 'ha os datadisk wipe')",
	Long:  osDataDiskWipeCmd.Long,
	Example: `
  ha factory-reset
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run:               osDataDiskWipeCmd.Run,
}

func init() {
	rootCmd.AddCommand(factoryResetCmd)
}

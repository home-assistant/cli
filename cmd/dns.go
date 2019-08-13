package cmd

import (
	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:     "dns",
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}

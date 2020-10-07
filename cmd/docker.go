package cmd

import (
	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:     "docker",
	Aliases: []string{"do"},
	Short:   "Docker backend specific for info and OCI configuration",
	Long: `
The docker command provides command-line tools to control the host docker that
Home Assistant is running on. It allows you to do things like use private OCI registries.`,
	Example: `
  ha docker info
  ha docker registries`,
}

func init() {
	rootCmd.AddCommand(dockerCmd)
}

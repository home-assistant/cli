package cmd

import (
	"github.com/spf13/cobra"
)

var jobsCmd = &cobra.Command{
	Use:     "jobs",
	Aliases: []string{"job", "tasks", "task"},
	Short:   "Get information and manage running jobs",
	Long: `
The jobs command allows you to manage the internal Home Assistant Job Manager and
exposing commands to view, configure and control it.`,
	Example: `
  ha jobs info
  ha jobs options`,
}

func init() {
	rootCmd.AddCommand(jobsCmd)
}

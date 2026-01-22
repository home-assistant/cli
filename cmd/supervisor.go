package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var supervisorCmd = &cobra.Command{
	Use:     "supervisor",
	Aliases: []string{"super", "su"},
	Short:   "Monitor, control and configure the Home Assistant Supervisor",
	Long: `
The Home Assistant Supervisor is the heart of the Home Assistant system.
It manages your Home Assistant Core, Operating System, and all the apps.
It even manages itself! This series of command give you control over the
Home Assistant Supervisor.`,
	Example: `
  ha supervisor reload
  ha supervisor update
  ha supervisor logs`,
}

func init() {
	log.Debug("Init supervisor")
	rootCmd.AddCommand(supervisorCmd)
}

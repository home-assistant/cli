package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var supervisorCmd = &cobra.Command{
	Use:     "supervisor",
	Aliases: []string{"super", "su"},
	Short:   "Monitor, control and configure the Hass.io Supervisor",
	Long: `
The Hass.io Supervisor is the core of the Hass.io system. It manages
your Home Assistant, HassOS, and all the add-ons. It even manages itself!
This series of command give you control over the Hass.io Supervisor.`,
	Example: `
  hassio supervisor reload
  hassio supervisor update
  hassio supervisor logs`,
}

func init() {
	log.Debug("Init supervisor")
	rootCmd.AddCommand(supervisorCmd)
}

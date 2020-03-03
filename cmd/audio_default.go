package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioDefaultCmd = &cobra.Command{
	Use:     "default",
	Aliases: []string{"def", "de", "standard", "std"},
	Short:   "Set default input/output audio device.",
	Long: `
Set the default input/output audio device of your Home Assistant system.
`,
	Example: `
	ha audio default input --name "..."
	ha audio default output --name "..."
`,
}

func init() {
	log.Debug("Init audio default")

	audioCmd.AddCommand(audioDefaultCmd)
}

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var audioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"snd", "sound"},
	Short:   "Audio device handling.",
	Long: `
Control audio devices.
`,
	Example: `
  ha audio info"
	`,
}

func init() {
	log.Debug("Init audio")

	rootCmd.AddCommand(audioCmd)
}

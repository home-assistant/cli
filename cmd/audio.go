package cmd

import (
	"log/slog"

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
  ha audio info
	`,
}

func init() {
	slog.Debug("Init audio")

	rootCmd.AddCommand(audioCmd)
}

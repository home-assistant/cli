package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:     "authentication",
	Aliases: []string{"auth", "au"},
	Short:   "Authentication for Home Assistant users.",
	Long: `
The authentication command allows you to manage Home Assistant user accounts.
`,
	Example: `
  ha authentication reset --username "JohnDoe" --password "123SuperSecret!"
	`,
}

func init() {
	log.Debug("Init authentication")

	rootCmd.AddCommand(authCmd)
}

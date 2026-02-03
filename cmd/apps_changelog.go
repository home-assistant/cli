package cmd

import (
	"errors"
	"fmt"
	"log/slog"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var appsChangelogCmd = &cobra.Command{
	Use:     "changelog [slug]",
	Aliases: []string{"cl", "ch"},
	Short:   "Show changelog of a Home Assistant app",
	Long: `
This command shows the changelog of an app. It gives you what has been
changed in the latest version and tell you about possible breaking changes.`,
	Example: `
ha apps changelog core_ssh
ha apps changelog core_mosquitto`,
	ValidArgsFunction: appsCompletions,
	Args:              cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("apps changelog", "args", args)

		section := "addons"
		command := "{slug}/changelog"

		url, err := helper.URLHelper(section, command)

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		request := helper.GetRequest()

		slug := args[0]

		request.SetPathParams(map[string]string{
			"slug": slug,
		})

		resp, err := request.SetHeader("Accept", "text/plain").Get(url)

		// returns 200 OK or 400, everything else is wrong
		if err == nil && resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			err = errors.New("unexpected server response")
			slog.Error("unexpected server response", "status", resp.StatusCode())
		}

		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			fmt.Println(string(resp.Body()))
		}
	},
}

func init() {
	appsCmd.AddCommand(appsChangelogCmd)
}

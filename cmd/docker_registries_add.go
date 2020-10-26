package cmd

import (
	"errors"
	"fmt"

	resty "github.com/go-resty/resty/v2"
	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dockerRegistriesAddCmd = &cobra.Command{
	Use:     "add [host]",
	Aliases: []string{"set", "new"},
	Short:   "Add new docker registry login for specific host",
	Long: `
Add new login for the Docker OCI registry server.
`,
	Example: `
  ha docker registries add my-docker.example.com --username "test" --password "example"
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("registries add")

		section := "docker"
		command := "registries"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}

		options := map[string]map[string]string{}

		request := helper.GetJSONRequest()

		host := args[0]
		options[host] = make(map[string]string)

		for _, value := range []string{
			"username",
			"password",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[host][value] = val
			}
		}

		if len(options) > 0 {
			log.WithField("options", options).Debug("Request body")
			request.SetBody(options)
		}

		resp, err := request.Post(url)

		// returns 200 OK or 400, everything else is wrong
		if err == nil {
			if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
				err = errors.New("Unexpected server response")
				log.Error(err)
			} else if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
				err = errors.New("API did not return a JSON response")
				log.Error(err)
			}
		}

		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	dockerRegistriesAddCmd.Flags().StringP("username", "u", "", "Username for OCI auth")
	dockerRegistriesAddCmd.Flags().StringP("password", "p", "", "Password for OCI auth")
	dockerRegistriesCmd.AddCommand(dockerRegistriesAddCmd)
}

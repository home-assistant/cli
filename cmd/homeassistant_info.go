package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	helper "github.com/home-assistant/hassio-cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	resty "gopkg.in/resty.v1"
	yaml "gopkg.in/yaml.v2"
)

// infoCmd represents the info command
var homeassistantInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"in"},
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("homeassistant info")

		section := "homeassistant"
		command := "info"
		base := viper.GetString("endpoint")

		url, err := helper.URLHelper(base, section, command)
		if err != nil {
			// TODO: error handler
			fmt.Printf("Error: %v", err)
			return
		}

		request := helper.GetClient()
		resp, err := request.Get(url)

		// explore response object
		log.WithFields(log.Fields{
			"statuscode":  resp.StatusCode(),
			"status":      resp.Status(),
			"time":        resp.Time(),
			"recieved-at": resp.ReceivedAt(),
			"headers":     resp.Header(),
			"request":     resp.Request.RawRequest,
			"body":        resp,
		}).Debug("Response")

		if !resty.IsJSONType(resp.Header().Get(http.CanonicalHeaderKey("Content-Type"))) {
			// TODO: return error
			fmt.Println("Error: api did not return a json response")
			return
		}

		var data map[string]interface{}

		if err := json.Unmarshal(resp.Body(), &data); err != nil {
			log.Fatal(err)
		} else {
			if data["result"] == "ok" {
				if len(data["data"].(map[string]interface{})) == 0 {
					fmt.Println("ok")
				} else {
					d, err := yaml.Marshal(data["data"])
					if err != nil {
						log.Fatalf("error: %v", err)
					}
					fmt.Print(string(d))
				}
			} else if data["result"] == "error" {
				os.Stderr.WriteString("ERROR\n")
				if data["message"] != nil {
					fmt.Fprintf(os.Stderr, "%v\n", data["message"])
				}
			} else {
				d, err := yaml.Marshal(data)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				fmt.Print(string(d))
			}
		}

		return
	},
}

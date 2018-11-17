package client

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	resty "gopkg.in/resty.v1"
	yaml "gopkg.in/yaml.v2"

	"net/url"
	"path"
	"strings"
)

var client *resty.Client

// URLHelper returns a url build from the arguments
func URLHelper(base, section, command string) (string, error) {
	log.WithFields(log.Fields{
		"base":    base,
		"section": section,
		"command": command,
	}).Debug("[GenerateURI]")

	scheme := ""
	if !strings.Contains(base, "://") {
		scheme = "http://"
	}

	uri := fmt.Sprintf("%s%s/%s/%s",
		scheme,
		base,
		section,
		command,
	)

	var url, err = url.Parse(uri)
	if err != nil {
		return "", err
	}

	url.Path = path.Clean(url.Path)
	log.WithFields(log.Fields{
		"uri": uri,
		"url": url.String(),
	}).Debug("[GenerateURI] Result")
	return url.String(), nil
}

// GetClient returns a resty.Request object prepared for a api call
func GetClient() *resty.Request {
	apiToken := viper.GetString("api-token")

	if client == nil {
		client = resty.New()
		// Registering Response Middleware
		client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
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

			return nil // if its success otherwise return error
		})
	}

	return client.R().
		SetHeader("Accept", "application/json").
		SetHeader("X-HASSIO-KEY", apiToken).
		SetAuthToken(apiToken)
}

// ShowJSONResponse formats a json response for human readers
func ShowJSONResponse(body []byte) {
	var data map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
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
}

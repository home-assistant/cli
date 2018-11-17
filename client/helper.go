package client

import (
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

// Response is the default json response fromt he home-assistant supervisor
type Response struct {
	Result  string                 `json:"result"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

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

// GetJSONRequest returns a request prepared for default json resposes
func GetJSONRequest() *resty.Request {
	request := GetRequest().
		SetResult(Response{})
	return request
}

// GetRequest returns a resty.Request object prepared for a api call
func GetRequest() *resty.Request {
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
func ShowJSONResponse(data *Response) {
	if data.Result == "ok" {
		if len(data.Data) == 0 {
			fmt.Println("ok")
		} else {
			d, err := yaml.Marshal(data.Data)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			fmt.Print(string(d))
		}
	} else if data.Result == "error" {
		os.Stderr.WriteString("ERROR\n")
		if data.Message != "" {
			fmt.Fprintf(os.Stderr, "%v\n", data.Message)
		}
	} else {
		d, err := yaml.Marshal(data)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Print(string(d))
	}
}

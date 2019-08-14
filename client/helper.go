package client

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	resty "github.com/go-resty/resty/v2"
	yaml "github.com/ghodss/yaml"

	"strings"
)

var client *resty.Client

// Response is the default JSON response from the Home Assistant Supervisor
type Response struct {
	Result  string                 `json:"result"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// URLHelper returns a URL build from the arguments
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

	var myurl, err = url.Parse(uri)
	if err != nil {
		return "", err
	}

	myurl.Path = path.Clean(myurl.Path)

	res, _ := url.PathUnescape(myurl.String())
	log.WithFields(log.Fields{
		"uri":         uri,
		"url":         myurl,
		"url(string)": res,
	}).Debug("[GenerateURI] Result")
	return res, nil
}

// GetJSONRequest returns a request prepared for default JSON resposes
func GetJSONRequest() *resty.Request {
	request := GetRequest().
		SetResult(Response{}).
		SetError(Response{})
	if RawJSON {
		request.
			SetDoNotParseResponse(true)
	}
	return request
}

// GetRequest returns a resty.Request object prepared for a API call
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

// ShowJSONResponse formats a JSON response for human readers
func ShowJSONResponse(resp *resty.Response) (success bool) {
	if RawJSON {
		// when we are returning raw JSON, all handling is for the consumer of the JSON
		success = true
		body := resp.RawBody()
		defer body.Close()
		if b, err := ioutil.ReadAll(body); err == nil {
			fmt.Print(string(b))
		}
		return
	}
	var data *Response
	if resp.IsSuccess() {
		data = resp.Result().(*Response)
	} else {
		data = resp.Error().(*Response)
	}
	if data.Result == "ok" {
		success = true
		if len(data.Data) == 0 {
			fmt.Println("ok")
		} else {
			j, err := json.Marshal(data.Data)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			d, err := yaml.JSONToYAML(j)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			fmt.Print(string(d))
		}
	} else if data.Result == "error" {
		fmt.Printf("Error: %s\n", data.Message)
	} else {
		d, err := yaml.Marshal(data)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Print(string(d))
	}
	return
}

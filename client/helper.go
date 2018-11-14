package client

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	resty "gopkg.in/resty.v1"

	"net/url"
	"path"
	"strings"
)

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
		"base":    base,
		"section": section,
		"command": command,
		"uri":     uri,
		"url":     url.String(),
	}).Debug("[GenerateURI] Result")
	return url.String(), nil
}

// GetClient returns a resty.Request object prepared for a api call
func GetClient() *resty.Request {
	apiToken := viper.GetString("api-token")

	return resty.R().
		SetHeader("Accept", "application/json").
		SetHeader("X-HASSIO-KEY", apiToken).
		SetAuthToken(apiToken)
}

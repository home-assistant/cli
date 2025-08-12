package client

import (
	"fmt"
	"time"

	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// RawJSON controls if the client does json handling or outputs it raw
var RawJSON = false

func GenericJSONErrorHandling(resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode() {
	case 200, 400, 404, 503:
		// Success and errors that should have JSON responses
		if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			return nil, fmt.Errorf("API did not return a JSON response. Status code %d", resp.StatusCode())
		}
	case 401:
		if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			return nil, fmt.Errorf("unauthorized: missing or invalid API token")
		}
	case 403:
		// Handle both JSON and plain text forbidden responses
		if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			return nil, fmt.Errorf("forbidden: insufficient permissions or invalid token")
		}
	case 502:
		return nil, fmt.Errorf("bad gateway: core proxy or ingress service failure")
	default:
		return nil, fmt.Errorf("unexpected server response. Status code: %d", resp.StatusCode())
	}

	return resp, err
}

func genericJSONMethod(get bool, section, command string, body map[string]any, timeout time.Duration) (*resty.Response, error) {
	url, err := URLHelper(section, command)
	if err != nil {
		return nil, err
	}

	request := GetJSONRequestTimeout(timeout)
	var resp *resty.Response

	if get {
		resp, err = request.Get(url)
	} else {
		if len(body) > 0 {
			log.WithField("body", body).Debug("Request body")
			request.SetBody(body)
		}
		resp, err = request.Post(url)
	}

	return GenericJSONErrorHandling(resp, err)
}

// GenericJSONGet is a helper for generic empty post request
func GenericJSONGet(section, command string) (*resty.Response, error) {
	return genericJSONMethod(true, section, command, nil, DefaultTimeout)
}

func GenericJSONGetTimeout(section, command string, timeout time.Duration) (*resty.Response, error) {
	return genericJSONMethod(true, section, command, nil, timeout)
}

// GenericJSONPost is a helper for generic empty post request
func GenericJSONPost(section, command string, body map[string]any) (*resty.Response, error) {
	return genericJSONMethod(false, section, command, body, DefaultTimeout)
}

func GenericJSONPostTimeout(section, command string, body map[string]any, timeout time.Duration) (*resty.Response, error) {
	return genericJSONMethod(false, section, command, body, timeout)
}

package client

import (
	"fmt"
	"net/http"
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

	// This is meant to handle HTTP response code as they are returned by Supervisor.
	// Handle known error codes as well as error codes which Supervisor returns without
	// a (JSON) body.
	switch resp.StatusCode() {
	case http.StatusOK, http.StatusBadRequest, http.StatusNotFound, http.StatusServiceUnavailable, http.StatusTooManyRequests:
		// Success and these errors should have JSON responses
		if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			return nil, fmt.Errorf("unexpected non-JSON response (status: %d)", resp.StatusCode())
		}
	case http.StatusInternalServerError:
		if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			return nil, fmt.Errorf("unknown error occurred, check supervisor logs with 'ha supervisor logs")
		}
	case http.StatusUnauthorized:
		if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			return nil, fmt.Errorf("unauthorized: missing or invalid API token")
		}
	case http.StatusForbidden:
		// Handle both JSON and plain text forbidden responses
		if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			return nil, fmt.Errorf("forbidden: insufficient permissions or invalid token")
		}
	case http.StatusBadGateway:
		return nil, fmt.Errorf("bad gateway: core proxy or ingress service failure")
	default:
		return nil, fmt.Errorf("unexpected server response (status: %d)", resp.StatusCode())
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

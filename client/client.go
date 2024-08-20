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
	if err == nil {
		switch resp.StatusCode() {
		case 200, 400, 403, 503:
			break
		default:
			err = fmt.Errorf("Unexpected server response. Status code: %d", resp.StatusCode())
			log.Error(err)
			return nil, err
		}

		if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			err = fmt.Errorf("API did not return a JSON response. Status code %d", resp.StatusCode())
			log.Error(err)
			return nil, err
		}
	}
	return resp, err
}

func genericJSONMethod(get bool, section, command string, body map[string]interface{}, timeout time.Duration) (*resty.Response, error) {
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
func GenericJSONPost(section, command string, body map[string]interface{}) (*resty.Response, error) {
	return genericJSONMethod(false, section, command, body, DefaultTimeout)
}

func GenericJSONPostTimeout(section, command string, body map[string]interface{}, timeout time.Duration) (*resty.Response, error) {
	return genericJSONMethod(false, section, command, body, timeout)
}

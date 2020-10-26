package client

import (
	"errors"

	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// RawJSON controls if the client does json handling or outputs it raw
var RawJSON = false

func genericJSONMethod(get bool, base, section, command string, body map[string]interface{}) (*resty.Response, error) {
	url, err := URLHelper(base, section, command)
	if err != nil {
		return nil, err
	}

	request := GetJSONRequest()
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

	// returns 200 OK or 400, everything else is wrong
	if err == nil {
		if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			err = errors.New("Unexpected server response")
			log.Error(err)
			return nil, err
		} else if !resty.IsJSONType(resp.Header().Get("Content-Type")) {
			err = errors.New("API did not return a JSON response")
			log.Error(err)
			return nil, err
		}
	}
	return resp, err
}

// GenericJSONGet is a helper for generic empty post request
func GenericJSONGet(base, section, command string) (*resty.Response, error) {
	return genericJSONMethod(true, base, section, command, nil)
}

// GenericJSONPost is a helper for generic empty post request
func GenericJSONPost(base, section, command string, body map[string]interface{}) (*resty.Response, error) {
	return genericJSONMethod(false, base, section, command, body)
}

package client

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

// GenericJSONPost is a helper for generic empty post request
func GenericJSONPost(base, section, command string, body map[string]interface{}) (*resty.Response, error) {
	url, err := URLHelper(base, section, command)
	if err != nil {
		return nil, err
	}

	request := GetJSONRequest()

	if body != nil && len(body) > 0 {
		log.WithField("body", body).Debug("Request body")
		request.SetBody(body)
	}

	resp, err := request.Post(url)

	// returns 200 OK or 400, everything else is wrong
	if err == nil {
		if resp.StatusCode() != 200 && resp.StatusCode() != 400 {
			err = errors.New("Unexpected server response")
			log.Error(err)
			return nil, err
		} else if !resty.IsJSONType(resp.Header().Get(http.CanonicalHeaderKey("Content-Type"))) {
			err = errors.New("api did not return a json response")
			log.Error(err)
			return nil, err
		}
	}
	return resp, nil
}

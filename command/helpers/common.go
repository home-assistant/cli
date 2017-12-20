package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HASSIO_SERVER uri to connect to hass.io with
const HASSIO_SERVER string = "http://hassio"

func GenerateUri(basepath string, endpoint string) string {
	var uri bytes.Buffer
	uri.WriteString(HASSIO_SERVER)
	uri.WriteString("/")
	uri.WriteString(basepath)
	uri.WriteString("/")
	uri.WriteString(endpoint)
	return uri.String()
}

func RestCall(basepath string, endpoint string, payload string) string {
	uri := GenerateUri(basepath, endpoint)
	var response *http.Response
	var err error

	if payload == "" {
		response, err = http.Get(uri)
	} else {
		jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
		jsonValue, _ := json.Marshal(jsonData)
		response, err = http.Post(uri, "application/json", bytes.NewBuffer(jsonValue))
	}

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		strData := string(data)
		return strData
	}
	return ""
}

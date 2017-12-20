package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// HASSIO_SERVER uri to connect to hass.io with
const HASSIO_SERVER = "http://hassio"

func GenerateUri(basepath string, endpoint string) string {
	var uri bytes.Buffer
	uri.WriteString(HASSIO_SERVER)
	uri.WriteString("/")
	uri.WriteString(basepath)
	uri.WriteString("/")
	uri.WriteString(endpoint)
	return uri.String()
}

func CreateJSONData(data string) map[string]string {
	var jsonData map[string]string
	var ss []string
	ss = strings.Split(data, ",")
	jsonData = make(map[string]string)
	for _, pair := range ss {
		z := strings.Split(pair, "=")
		jsonData[z[0]] = z[1]
	}
	return jsonData
}

func RestCall(basepath string, endpoint string, get bool, payload string) map[string]interface{} {
	uri := GenerateUri(basepath, endpoint)
	var response *http.Response
	var err error

	if get {
		response, err = http.Get(uri)
	} else {
		jsonData := CreateJSONData(payload)
		jsonValue, _ := json.Marshal(jsonData)
		response, err = http.Post(uri, "application/json", bytes.NewBuffer(jsonValue))
	}

	if err != nil {
		fmt.Printf("The HTTP request failed with the error: %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var f interface{}
		err := json.Unmarshal(data, &f)
		if err != nil {
			fmt.Printf("Error decoding json %s", err)
			return nil
		}
		res := f.(map[string]interface{})
		return res
	}
	return nil
}

func DisplayOutput(data string, json bool) {
	if !(json) {
		// Make data pretty
		//fmt.Println(data)
	}
	fmt.Println(data)
}

func MapToJSON(data map[string]interface{}) string {
	b, _ := json.Marshal(data)
	// Convert bytes to string.
	s := string(b)
	return s
}

func FilterProperties(data string, filter []string) []string {
	return nil
}
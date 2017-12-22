package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"os"
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

func RestCall(basepath string, endpoint string, bGet bool, payload string) string {
	uri := GenerateUri(basepath, endpoint)
	var response *http.Response
	var err error

	if bGet {
		response, err = http.Get(uri)
	} else {
		jsonData := CreateJSONData(payload)
		jsonValue, _ := json.Marshal(jsonData)
		response, err = http.Post(uri, "application/json", bytes.NewBuffer(jsonValue))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "The HTTP request failed with the error: %s\n", err)
		os.Exit(1)
	}
	data, _ := ioutil.ReadAll(response.Body)
	return string(data)
}

func StrToMap(data string) map[string]interface{} {
	var b = []byte(data)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding json %s: %s", err, data)
		os.Exit(4)
	}
	res := f.(map[string]interface{})
	return res
}

func DisplayOutput(data string, rawjson bool) {
	if rawjson {
		fmt.Println(data)
	} else {
		b := []byte(data)
		x := bytes.Buffer{}
		json.Indent(&x, b, "", "    ")
		fmt.Println(string(x.Bytes()))
	}
}

func MapToJSON(data map[string]interface{}) string {
	b, _ := json.Marshal(data)
	// Convert bytes to string.
	s := string(b)
	return s
}

func FilterProperties(data string, filter []string) map[string]string {
	mymap := StrToMap(data)
	newmap := make(map[string]string)
	for _, value := range filter {
		if val, ok := mymap[value]; ok {
			newmap[value] = val.(string)
		}
	}
	return newmap
}


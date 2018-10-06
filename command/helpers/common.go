package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// HassioServer uri to connect to hass.io with
const HassioServer = "hassio"

// GenerateURI Creates the API URI from the server and the endpoint
func GenerateURI(basepath string, endpoint string, serverOverride string) string {
	log.WithFields(log.Fields{
		"basepath":       basepath,
		"endpoint":       endpoint,
		"serverOverride": serverOverride,
	}).Debug("[GenerateURI]")
	var uri bytes.Buffer
	uri.WriteString("http://")
	if serverOverride != "" {
		uri.WriteString(serverOverride)
	} else if os.Getenv("HASSIO") != "" {
		uri.WriteString(os.Getenv("HASSIO"))
	} else {
		uri.WriteString(HassioServer)
	}
	uri.WriteString("/")
	uri.WriteString(basepath)
	if endpoint != "" {
		uri.WriteString("/")
		uri.WriteString(endpoint)
	}
	return uri.String()
}

func CreateJSONData(data string) map[string]string {
	log.WithFields(log.Fields{
		"data": data,
	}).Debug("[CreateJSONData]")

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

// RestCall Makes the Call to the API
func RestCall(uri string, bGet bool, payload string) []byte {
	var response *http.Response
	var request *http.Request
	var err error
	var client = &http.Client{}
	var XHassioKey = os.Getenv("HASSIO_TOKEN")

	log.WithFields(log.Fields{
		"uri":     uri,
		"bGet":    bGet,
		"payload": payload,
	}).Debug("[RestCall]")

	if bGet {
		request, err = http.NewRequest("GET", uri, nil)
	} else {
		if payload != "" {
			jsonValue := []byte("")
			jsonData := CreateJSONData(payload)
			jsonValue, _ = json.Marshal(jsonData)

			request, err = http.NewRequest("POST", uri, bytes.NewBuffer(jsonValue))
			request.Header.Add("contentType", "application/json")
		} else {
			request, err = http.NewRequest("POST", uri, nil)
		}
	}

	request.Header.Add("X-HASSIO-KEY", XHassioKey)
	response, err = client.Do(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "The HTTP request failed with the error: %s\n", err)
		os.Exit(1)
	}
	data, _ := ioutil.ReadAll(response.Body)
	log.WithFields(log.Fields{
		"data": data,
	}).Debug("[RestCall]")

	defer response.Body.Close()
	return data
}

func ByteArrayToMap(data []byte) map[string]interface{} {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding json %s: %s\n", err, string(data))
		os.Exit(4)
	}
	res := f.(map[string]interface{})
	return res
}

func DisplayOutput(data []byte, rawjson bool) {
	if rawjson {
		fmt.Println(string(data))
	} else {
		mymap := ByteArrayToMap(data)
		if mymap["result"] == "ok" && len(mymap["data"].(map[string]interface{})) == 0 {
			fmt.Println(mymap["result"])
		} else if mymap["result"] == "error" {
			os.Stderr.WriteString("ERROR\n")
			if mymap["message"] != nil {
				fmt.Fprintf(os.Stderr, "%v\n", mymap["message"])
			}
		} else {
			x := bytes.Buffer{}
			json.Indent(&x, data, "", "    ")
			fmt.Println(string(x.Bytes()))
		}
	}
}

func FilterProperties(data []byte, filter []string) []byte {
	log.WithFields(log.Fields{
		"indata": string(data),
		"filter": filter,
	}).Debug("[FilterProperties]")

	mymap := ByteArrayToMap(data)
	mymapdata := mymap["data"].(map[string]interface{})
	newmap := make(map[string]interface{})
	for _, value := range filter {
		if val, ok := mymapdata[value]; ok {
			newmap[value] = fmt.Sprintf("%v", val)
		}
	}
	rawjson, _ := json.Marshal(newmap)
	log.WithFields(log.Fields{
		"outdata": string(rawjson),
	}).Debug("[FilterProperties]")

	return rawjson
}

// ExecCommand Used to execute the remote calls for each of the managing commands
func ExecCommand(basepath string, endpoint string, serverOverride string, get bool, Options string, Filter string, RawJSON bool) {
	log.WithFields(log.Fields{
		"basepath":       basepath,
		"endpoint":       endpoint,
		"serverOverride": serverOverride,
		"get":            get,
		"options":        Options,
		"rawJSON":        RawJSON,
		"filter":         Filter,
	}).Debug("[ExecCommand]")
	uri := GenerateURI(basepath, endpoint, serverOverride)
	response := RestCall(uri, get, Options)
	if Filter == "" {
		DisplayOutput(response, RawJSON)
	} else {
		filter := strings.Split(Filter, ",")
		data := FilterProperties(response, filter)
		DisplayOutput(data, RawJSON)
	}
	responseMap := ByteArrayToMap(response)
	if responseMap["result"] != "ok" {
		os.Exit(10)
	}
}

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

// HassioServer uri to connect to hass.io with
const HassioServer = "http://hassio"

// GenerateUri Creates the API URI from the server and the endpoint
func GenerateUri(basepath string, endpoint string, serverOverride string) string {
    var uri bytes.Buffer
    if serverOverride != "" {
        uri.WriteString(serverOverride)
    } else {
        uri.WriteString(HassioServer)
    }
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

// RestCall Makes the Call to the API
func RestCall(uri string, bGet bool, payload string) []byte {
    var response *http.Response
    var request *http.Request
    var err error
    var client = &http.Client{}

    if bGet {
        request, err = http.NewRequest("GET", uri, nil)
        request.Header.Add("X_HASSIO_KEY", os.Getenv("X-HASSIO-KEY"))
    } else {
        jsonData := CreateJSONData(payload)
        jsonValue, _ := json.Marshal(jsonData)

        request, err = http.NewRequest("POST", uri, bytes.NewBuffer(jsonValue))
        request.Header.Add("X-HASSIO-KEY", os.Getenv("X-HASSIO-KEY"))
        request.Header.Add("contentType", "application/json")
    }

    response, err = client.Do(request)
    defer response.Body.Close()

    if err != nil {
        fmt.Fprintf(os.Stderr, "The HTTP request failed with the error: %s\n", err)
        os.Exit(1)
    }
    data, _ := ioutil.ReadAll(response.Body)
    return data
}

func ByteArrayToMap(data []byte) map[string]interface{} {
    var f interface{}
    err := json.Unmarshal(data, &f)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error decoding json %s: %s", err, string(data))
        os.Exit(4)
    }
    res := f.(map[string]interface{})
    return res
}

func DisplayOutput(data []byte, rawjson bool) {
    if rawjson {
        fmt.Println(string(data))
    } else {
        x := bytes.Buffer{}
        json.Indent(&x, data, "", "    ")
        fmt.Println(string(x.Bytes()))
    }
}

func FilterProperties(data []byte, filter []string) []byte {
    mymap := ByteArrayToMap(data)
    mymapdata := mymap["data"].(map[string]interface{})
    newmap := make(map[string]interface{})
    for _, value := range filter {
        if val, ok := mymapdata[value]; ok {
            newmap[value] = fmt.Sprintf("%v",val)
        }
    }
    rawjson, _ := json.Marshal(newmap)
    return rawjson
}
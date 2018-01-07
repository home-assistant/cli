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
// DebugEnabled set by global --debug (-d) flag
var DebugEnabled = false

// GenerateURI Creates the API URI from the server and the endpoint
func GenerateURI(basepath string, endpoint string, serverOverride string) string {
    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [GenerateURI]: basepath->'%s', endpoint->'%s', serverOverride->'%s'\n", basepath, endpoint, serverOverride)
    }
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
    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [CreateJSONData]: data->'%s'\n", data)
    }
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

    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [RestCall]: data->'%s', GET->'%t', payload->'%s'\n", uri, bGet, payload)
    }

    if bGet {
        request, err = http.NewRequest("GET", uri, nil)
        request.Header.Add("X_HASSIO_KEY", os.Getenv("X-HASSIO-KEY"))
    } else {
        jsonValue := []byte("")
        if payload != "" {
            jsonData := CreateJSONData(payload)
            jsonValue, _ = json.Marshal(jsonData)
        }

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
    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [RestCall]: ResponseBody->'%s'\n", string(data))
    }
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
    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [FilterProperties]: indata->'%s', filter->'%s'\n", string(data), filter)
    }
    mymap := ByteArrayToMap(data)
    mymapdata := mymap["data"].(map[string]interface{})
    newmap := make(map[string]interface{})
    for _, value := range filter {
        if val, ok := mymapdata[value]; ok {
            newmap[value] = fmt.Sprintf("%v",val)
        }
    }
    rawjson, _ := json.Marshal(newmap)
    if DebugEnabled {
        fmt.Fprintf(os.Stdout, "DEBUG [FilterProperties]: outdata->'%s'\n", string(rawjson))
    }
    return rawjson
}

// ExecCommand Used to execute the remote calls for each of the managing commands
func ExecCommand(basepath string, endpoint string, serverOverride string, get bool, Options string, Filter string, RawJSON bool) {
    uri := GenerateURI(basepath, endpoint, serverOverride)
    response := RestCall(uri, get,  Options)
    if Filter == "" {
        DisplayOutput(response, RawJSON)
    } else {
        filter := strings.Split(Filter, ",")
        data := FilterProperties(response, filter)
        DisplayOutput(data, RawJSON)
    }
    responseMap := ByteArrayToMap(response)
    result := responseMap["result"]
    if result != "ok" {
        os.Exit(10)
    }
}
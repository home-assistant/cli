package helpers

import (
    "testing"
    "bytes"
)

func TestRestCall(t *testing.T) {
    // Write your code here
}

func TestGenerateUri(t *testing.T) {
    expectedURI := "http://hassio/api/endpoint"
    uri := GenerateURI("api", "endpoint", "")
    if uri != expectedURI {
        t.Errorf("URI incorrect, got: %s, want: %s.", uri, expectedURI)
    }
}

func TestGenerateUriOverride(t *testing.T) {
    expectedURI := "http://testme/api/endpoint"
    uri := GenerateURI("api", "endpoint", "http://testme")
    if uri != expectedURI {
        t.Errorf("URI incorrect, got: %s, want: %s.", uri, expectedURI)
    }
}

func TestCreateJSONData(t *testing.T) {
    expectedVersion := "0.23"
    res := CreateJSONData("version=0.23")
    if res["version"] != expectedVersion {
        t.Errorf("Value mismatch, got: %s, want: %s.", res["version"], expectedVersion)
    }
}

func TestCreateJSONData_multi(t *testing.T) {
    expectedVersion := "0.23"
    expectedOther := "yes"
    res := CreateJSONData("version=0.23,other=yes")
    if res["version"] != expectedVersion {
        t.Errorf("Value mismatch, got: %s, want: %s.", res["version"], expectedVersion)
    }
    if res["other"] != expectedOther {
        t.Errorf("Value mismatch, got: %s, want: %s.", res["other"], expectedOther)
    }
}

func TestFilterProperties(t *testing.T) {
    //supervisor info
    expected := `{"version":"0.60"}`
    rawdata := []byte(`{"result": "ok", "data": {"version": "0.60", "last_version": "0.60", "image": "homeassistant/raspberrypi2-homeassistant", "devices": [], "custom": false, "boot": true, "port": 8123, "ssl": false, "watchdog": true}}`)
    filter := []string{"version"}
    res := FilterProperties(rawdata, filter)
    if string(res) != expected {
        t.Errorf("Value mismatch, got: %s, want: %s.", res, expected)
    }
}

func TestByteArrayToMap(t *testing.T) {
    expectedStr := "TestVal"
    var myStr bytes.Buffer
    myStr.WriteString(`{"TestKey":"TestVal"}`)
    res := ByteArrayToMap(myStr.Bytes())

    if res["TestKey"] != expectedStr {
        t.Errorf("Value mismatch, got: %s, want: %s.", res["TestKey"], expectedStr)
    }
}

// host hardware
//rawdata := []byte(`{"result": "ok", "data": {"serial": ["/dev/ttyACM0", "/dev/ttyAMA0"], "input": [], "disk": [], "gpio": ["gpiochip0"], "audio": {"0": {"name": "bcm2835 - bcm2835 ALSA", "type": "ALSA", "devices": {"0": "digital audio playback", "1": "digital audio playback"}}}}}`)
// homeassistant info
//rawdata := []byte(`{"result": "ok", "data": {"version": "0.60", "last_version": "0.60", "image": "homeassistant/raspberrypi2-homeassistant", "devices": [], "custom": false, "boot": true, "port": 8123, "ssl": false, "watchdog": true}}`)

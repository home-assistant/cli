package helpers

import (
	"testing"
)

func TestRestCall(t *testing.T) {
	// Write your code here
}

func TestGenerateUri(t *testing.T) {
	expectedURI := "http://hassio/api/endpoint"
	uri := GenerateUri("api", "endpoint")
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
	expectedversion := "0.76"
	rawdata := string([]byte(`{"version": "0.76", "last_version": "0.76", "beta_channel": false, "arch": "armhf", "timezone": "Europe/London", "addons": [{"name": "Git pull", "slug": "core_git_pull", "description": "Simple git pull to update the local configuration", "state": "started", "version": "2.3", "installed": "2.2", "repository": "core", "logo": true}, {"name": "Check Home Assistant configuration", "slug": "core_check_config", "description": "Check current Home Assistant configuration against a new version", "state": "stopped", "version": "0.7", "installed": "0.7", "repository": "core", "logo": true}, {"name": "SSH - Secure Shell", "slug": "a0d7b954_ssh", "description": "Allows SSH connections to your Home Assistant instance", "state": "started", "version": "2.1.0", "installed": "2.1.0", "repository": "a0d7b954", "logo": true}, {"name": "IDE", "slug": "a0d7b954_ide", "description": "Advanced IDE for Home Assistant, based on Cloud9 IDE", "state": "started", "version": "0.1.0", "installed": "0.1.0", "repository": "a0d7b954", "logo": false}, {"name": "Terminal", "slug": "a0d7b954_terminal", "description": "Terminal access to your Home Assistant instance via the web", "state": "started", "version": "2.1.0", "installed": "2.1.0", "repository": "a0d7b954", "logo": true}], "addons_repositories": ["https://github.com/hassio-addons/repository"]}`))
	filter := []string{"version"}
	// host hardware
	//rawdata := []byte(`{"result": "ok", "data": {"serial": ["/dev/ttyACM0", "/dev/ttyAMA0"], "input": [], "disk": [], "gpio": ["gpiochip0"], "audio": {"0": {"name": "bcm2835 - bcm2835 ALSA", "type": "ALSA", "devices": {"0": "digital audio playback", "1": "digital audio playback"}}}}}`)
	// homeassistant info
	//rawdata := []byte(`{"result": "ok", "data": {"version": "0.60", "last_version": "0.60", "image": "homeassistant/raspberrypi2-homeassistant", "devices": [], "custom": false, "boot": true, "port": 8123, "ssl": false, "watchdog": true}}`)
	res := FilterProperties(rawdata, filter)
	if res["version"] != expectedversion {
		t.Errorf("Value mismatch, got: %s, want: %s.", res["version"], expectedversion)
	}
}

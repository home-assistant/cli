package helpers

import "testing"

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

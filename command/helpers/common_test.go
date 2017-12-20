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

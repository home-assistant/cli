package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRequestRejectsHTTPS(t *testing.T) {
	_, err := GetJSONRequest().Get("https://supervisor/supervisor/info")
	if err == nil || !strings.Contains(err.Error(), "unsupported protocol scheme") {
		t.Errorf("got error %v", err)
	}
}

func TestRequestConnectionRefused(t *testing.T) {
	// Get a port nothing listens on
	server := httptest.NewServer(http.NotFoundHandler())
	url := server.URL
	server.Close()

	_, err := GetJSONRequest().Get(url + "/supervisor/info")
	if err == nil || !strings.HasPrefix(err.Error(), `Get "`+url+`/supervisor/info": dial tcp`) {
		t.Errorf("got error %v", err)
	}
}

func TestRequestCloseEndlessStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "first\n")
		w.(http.Flusher).Flush()
		// Keep the response open until the client disconnects, like
		// following logs does
		<-r.Context().Done()
	}))
	t.Cleanup(server.Close)

	resp, err := GetRequestTimeout(0).SetDoNotParseResponse(true).Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	line := make([]byte, len("first\n"))
	if _, err := io.ReadFull(resp.RawBody(), line); err != nil || string(line) != "first\n" {
		t.Errorf("got %q, error %v", line, err)
	}

	// Closing must not wait for the end of the response
	closed := make(chan error, 1)
	go func() { closed <- resp.RawBody().Close() }()
	select {
	case err := <-closed:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(5 * time.Second):
		t.Error("closing the body blocks")
	}
}

func TestRequestDoesNotFollowRedirects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/elsewhere", http.StatusFound)
	}))
	t.Cleanup(server.Close)

	resp, err := GetRequest().Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode() != http.StatusFound {
		t.Errorf("got status %d, want %d", resp.StatusCode(), http.StatusFound)
	}
}

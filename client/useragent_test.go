package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestUserAgent(t *testing.T) {
	old := Version
	t.Cleanup(func() { Version = old })

	Version = "5.6.0"
	if got, want := userAgent(), "HomeAssistantCLI/5.6.0"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	// Test binaries carry no module version
	Version = ""
	if got, want := userAgent(), "HomeAssistantCLI/dev"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestUserAgentSent(t *testing.T) {
	var got string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok","data":{}}`))
	}))
	t.Cleanup(server.Close)

	viper.SetDefault("endpoint", server.Listener.Addr().String())
	client = nil
	t.Cleanup(func() { client = nil })

	if _, err := GenericJSONGet("supervisor", "info"); err != nil {
		t.Fatal(err)
	}
	if want := userAgent(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

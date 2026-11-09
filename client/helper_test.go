package client

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

var urlTests = []struct {
	base    string
	section string
	command string
	out     string
}{
	{"supervisor", "section", "command", "http://supervisor/section/command"},
	{"supervisor:80", "section", "command", "http://supervisor:80/section/command"},
	{"supervisor.example.org:8080", "section", "command", "http://supervisor.example.org:8080/section/command"},
	{"https://supervisor", "section", "command", "https://supervisor/section/command"},
	{"https://supervisor:8080", "section", "command", "https://supervisor:8080/section/command"},
	{"https://supervisor", "section", "command", "https://supervisor/section/command"},
	{"https://supervisor:8080", "section", "command", "https://supervisor:8080/section/command"},
	{"supervisor", "section", "", "http://supervisor/section"},
	{"supervisor", "section/../othersection", "", "http://supervisor/othersection"},
	{"supervisor/api/", "section", "command", "http://supervisor/api/section/command"},
	{"supervisor/api/", "section", "{slug}/command", "http://supervisor/api/section/{slug}/command"},
	{"unix:///run/supervisor/api/supervisor.sock", "section", "command", "http://localhost/section/command"},
	{"unix:///run/supervisor/api/supervisor.sock", "section", "", "http://localhost/section"},
}

func TestURLHelper(t *testing.T) {
	for _, tt := range urlTests {
		t.Run(fmt.Sprintf("[%s][%s][%s]", tt.base, tt.section, tt.command), func(t *testing.T) {
			viper.SetDefault("endpoint", tt.base)
			s, _ := URLHelper(tt.section, tt.command)
			if s != tt.out {
				t.Errorf("got %q, want %q", s, tt.out)
			}
		})
	}
}

func TestURLHelperUnixSocketWithoutPath(t *testing.T) {
	viper.SetDefault("endpoint", "unix://")
	if _, err := URLHelper("section", "command"); err == nil {
		t.Error("expected an error for a unix:// endpoint without a path")
	}
}

func TestGetRequest(t *testing.T) {
	tests := []struct {
		name      string
		overUnix  bool
		token     string
		wantToken string
	}{
		{"tcp with token", false, "abc", "Bearer abc"},
		{"unix socket without token", true, "", ""},
		{"unix socket ignores token", true, "abc", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotAuth, gotHost, gotPath string
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotAuth = r.Header.Get("Authorization")
				gotHost = r.Host
				gotPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"result":"ok","data":{}}`))
			}))
			t.Cleanup(server.Close)

			var endpoint string
			if tt.overUnix {
				// Keep the path short, sun_path is limited to ~107 characters
				dir, err := os.MkdirTemp("", "ha")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { os.RemoveAll(dir) })
				socket := filepath.Join(dir, "supervisor.sock")
				listener, err := net.Listen("unix", socket)
				if err != nil {
					t.Fatal(err)
				}
				// Replace the TCP listener the unstarted server already opened
				server.Listener.Close()
				server.Listener = listener
				endpoint = UnixSocketScheme + socket
			} else {
				endpoint = server.Listener.Addr().String()
			}
			server.Start()

			viper.SetDefault("endpoint", endpoint)
			viper.SetDefault("api-token", tt.token)
			client = nil
			t.Cleanup(func() { client = nil })

			resp, err := GenericJSONGet("supervisor", "info")
			if err != nil {
				t.Fatal(err)
			}
			if !resp.IsSuccess() {
				t.Fatalf("unexpected status %d", resp.StatusCode())
			}
			if gotPath != "/supervisor/info" {
				t.Errorf("path: got %q, want %q", gotPath, "/supervisor/info")
			}
			if gotAuth != tt.wantToken {
				t.Errorf("Authorization: got %q, want %q", gotAuth, tt.wantToken)
			}
			if tt.overUnix && gotHost != "localhost" {
				t.Errorf("Host: got %q, want %q", gotHost, "localhost")
			}
		})
	}
}

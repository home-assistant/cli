package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	resty "github.com/go-resty/resty/v2"
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

// resetClients drops cached clients so a test starts from a fresh state.
func resetClients(t *testing.T) {
	t.Helper()
	reset := func() {
		clients = map[time.Duration]*resty.Client{}
		transport = nil
	}
	reset()
	t.Cleanup(reset)
}

func TestGetRequestTimeoutPerRequest(t *testing.T) {
	const delay = 200 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/slow" {
			time.Sleep(delay)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"ok","data":{}}`))
	}))
	t.Cleanup(server.Close)
	resetClients(t)

	url := server.URL

	// A first request with a short timeout must not cap later ones
	if _, err := GetRequestTimeout(delay / 4).Get(url + "/fast"); err != nil {
		t.Fatal(err)
	}
	if _, err := GetRequestTimeout(10 * delay).Get(url + "/slow"); err != nil {
		t.Errorf("longer timeout not applied: %v", err)
	}
	if _, err := GetRequestTimeout(delay / 4).Get(url + "/slow"); err == nil {
		t.Error("shorter timeout not applied")
	}

	if len(clients) != 2 {
		t.Errorf("got %d clients, want 2", len(clients))
	}
	for timeout, c := range clients {
		if c.GetClient().Transport != transport {
			t.Errorf("client for %s does not use the shared transport", timeout)
		}
	}
}

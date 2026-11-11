package client

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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

// captureStdout returns what f writes to stdout.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	out := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		out <- string(b)
	}()

	f()
	_ = w.Close()
	return <-out
}

func TestShowJSONResponseRawJSON(t *testing.T) {
	RawJSON = true
	defer func() { RawJSON = false }()

	tests := []struct {
		name    string
		status  int
		body    string
		success bool
	}{
		{"ok", http.StatusOK, `{"result":"ok","data":{"version":"1.0"}}`, true},
		{"ok without data", http.StatusOK, `{"result":"ok"}`, true},
		{"ok with non-object data", http.StatusOK, `{"result":"ok","data":["a","b"]}`, true},
		{"error", http.StatusBadRequest, `{"result":"error","message":"App is already installed"}`, false},
		{"error with success status", http.StatusOK, `{"result":"error","message":"failed"}`, false},
		{"unknown result", http.StatusOK, `{"result":"unknown"}`, false},
		{"missing result", http.StatusOK, `{}`, false},
		{"empty body", http.StatusOK, ``, false},
		{"not JSON", http.StatusBadGateway, `<html>Bad Gateway</html>`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = io.WriteString(w, tt.body)
			}))
			defer server.Close()

			resp, err := GetJSONRequest().Get(server.URL)
			if err != nil {
				t.Fatal(err)
			}

			var success bool
			output := captureStdout(t, func() {
				success = ShowJSONResponse(resp)
			})

			if success != tt.success {
				t.Errorf("success = %v, want %v", success, tt.success)
			}
			if output != tt.body {
				t.Errorf("output = %q, want the body as-is %q", output, tt.body)
			}
		})
	}
}

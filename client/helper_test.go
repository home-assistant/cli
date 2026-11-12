package client

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	{"unix", "section", "command", "http://unix/section/command"},
	{"unix:8080", "section", "command", "http://unix:8080/section/command"},
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

func TestURLHelperInvalidUnixSocket(t *testing.T) {
	for _, endpoint := range []string{"unix://", "unix:/run/x.sock"} {
		t.Run(endpoint, func(t *testing.T) {
			viper.SetDefault("endpoint", endpoint)
			if _, err := URLHelper("section", "command"); err == nil {
				t.Errorf("expected an error for endpoint %q", endpoint)
			}
		})
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
	viper.SetDefault("endpoint", server.Listener.Addr().String())
	client = nil
	t.Cleanup(func() { client = nil })

	url := server.URL

	// A first request with a short timeout must not cap later ones
	if _, err := GetRequestTimeout(delay / 4).Get(url + "/fast"); err != nil {
		t.Fatal(err)
	}
	first := client
	if _, err := GetRequestTimeout(10 * delay).Get(url + "/slow"); err != nil {
		t.Errorf("longer timeout not applied: %v", err)
	}
	if _, err := GetRequestTimeout(delay / 4).Get(url + "/slow"); err == nil {
		t.Error("shorter timeout not applied")
	}
	if _, err := GetRequestTimeout(0).Get(url + "/slow"); err != nil {
		t.Errorf("timeout applied although disabled: %v", err)
	}

	if client != first {
		t.Error("requests do not share one client")
	}
}

func TestRequest(t *testing.T) {
	var gotMethod, gotPath, gotContentType, gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.EscapedPath()
		gotContentType = r.Header.Get("Content-Type")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"result":"error","message":"nope"}`)
			return
		}
		_, _ = io.WriteString(w, `{"result":"ok","data":{"key":"value"}}`)
	}))
	t.Cleanup(server.Close)
	viper.SetDefault("endpoint", server.Listener.Addr().String())
	client = nil
	t.Cleanup(func() { client = nil })

	t.Run("post with body and path params", func(t *testing.T) {
		resp, err := GetJSONRequest().
			SetPathParams(map[string]string{"slug": "a/b c"}).
			SetBody(map[string]any{"version": "1.0"}).
			Post(server.URL + "/addons/{slug}/update")
		if err != nil {
			t.Fatal(err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("method: got %q", gotMethod)
		}
		// The parameter is escaped, it cannot add path segments
		if want := "/addons/a%2Fb%20c/update"; gotPath != want {
			t.Errorf("path: got %q, want %q", gotPath, want)
		}
		if gotContentType != "application/json" {
			t.Errorf("Content-Type: got %q", gotContentType)
		}
		if want := `{"version":"1.0"}`; gotBody != want {
			t.Errorf("body: got %q, want %q", gotBody, want)
		}
		if !resp.IsSuccess() || resp.Result().Result != "ok" || resp.Result().Data["key"] != "value" {
			t.Errorf("unexpected result %+v", resp.Result())
		}
	})

	t.Run("get without body", func(t *testing.T) {
		if _, err := GetJSONRequest().Get(server.URL + "/info"); err != nil {
			t.Fatal(err)
		}
		if gotMethod != http.MethodGet || gotBody != "" || gotContentType != "" {
			t.Errorf("got method %q, body %q, Content-Type %q", gotMethod, gotBody, gotContentType)
		}
	})

	t.Run("error response", func(t *testing.T) {
		resp, err := GetJSONRequest().Delete(server.URL + "/mounts/x")
		if err != nil {
			t.Fatal(err)
		}
		if resp.IsSuccess() || resp.StatusCode() != http.StatusBadRequest {
			t.Errorf("unexpected status %d", resp.StatusCode())
		}
		if resp.Error().Result != "error" || resp.Error().Message != "nope" {
			t.Errorf("unexpected error %+v", resp.Error())
		}
		if resp.Result().Result != "" {
			t.Errorf("result set for an error response: %+v", resp.Result())
		}
	})

	t.Run("stream", func(t *testing.T) {
		request := GetRequestTimeout(0).SetDoNotParseResponse(true)
		request.URL = server.URL + "/logs"
		resp, err := request.Send()
		if err != nil {
			t.Fatal(err)
		}
		var success bool
		output := captureStdout(t, func() {
			success = StreamTextResponse(resp)
		})
		if want := `{"result":"ok","data":{"key":"value"}}`; !success || output != want {
			t.Errorf("got success %v, output %q", success, output)
		}
	})
}

func TestRequestTimeoutStreamedBody(t *testing.T) {
	const delay = 200 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send the headers right away, the body only after the delay
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		select {
		case <-time.After(delay):
		case <-r.Context().Done():
		}
		_, _ = io.WriteString(w, "late")
	}))
	t.Cleanup(server.Close)
	viper.SetDefault("endpoint", server.Listener.Addr().String())
	client = nil
	t.Cleanup(func() { client = nil })

	// The timeout also covers reading a body the caller streams
	resp, err := GetRequestTimeout(delay / 4).SetDoNotParseResponse(true).Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(resp.RawBody()); err == nil {
		t.Error("timeout not applied while reading the body")
	}
	_ = resp.RawBody().Close()

	resp, err = GetRequestTimeout(10 * delay).SetDoNotParseResponse(true).Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(resp.RawBody())
	_ = resp.RawBody().Close()
	if err != nil || string(b) != "late" {
		t.Errorf("got body %q, error %v", b, err)
	}
}

func TestIsJSONType(t *testing.T) {
	for contentType, want := range map[string]bool{
		"application/json":                true,
		"application/json; charset=utf-8": true,
		"application/problem+json":        true,
		"text/json":                       true,
		"Application/JSON":                true,
		"application/json; charset":       true,
		"text/plain":                      false,
		"text/html; charset=utf-8":        false,
		"":                                false,
	} {
		if got := IsJSONType(contentType); got != want {
			t.Errorf("IsJSONType(%q) = %v, want %v", contentType, got, want)
		}
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

			// The client is created for the configured endpoint, do not reuse
			// one a previous test may have set up for a Unix socket
			viper.SetDefault("endpoint", server.Listener.Addr().String())
			client = nil
			t.Cleanup(func() { client = nil })

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

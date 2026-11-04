package client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/term"

	yaml "github.com/ghodss/yaml"
	resty "github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"strings"
)

const DefaultTimeout = 30 * time.Second
const ContainerOperationTimeout = 10 * time.Minute
const ContainerDownloadTimeout = 1 * time.Hour
const OsDownloadTimeout = 1 * time.Hour
const BackupTimeout = 3 * time.Hour
const RebootTimeout = 90 * time.Second
const MountTimeout = 5 * time.Minute

// UnixSocketScheme is the endpoint prefix selecting a Unix socket, e.g.
// unix:///run/supervisor/api/supervisor.sock
const UnixSocketScheme = "unix://"

// unixSocketBaseURL is the base URL of requests sent over a Unix socket. The
// host is only used for the Host header, the connection goes to the socket.
const unixSocketBaseURL = "http://localhost"

var client *resty.Client

// Response is the default JSON response from the Home Assistant Supervisor
type Response struct {
	Result  string         `json:"result"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

// socketPath returns the Unix socket path of a unix:// endpoint. ok is false
// for any other endpoint.
func socketPath(endpoint string) (path string, ok bool) {
	path, ok = strings.CutPrefix(endpoint, UnixSocketScheme)
	if !ok {
		return "", false
	}
	return path, true
}

// URLHelper returns a URL built from the arguments
func URLHelper(section, command string) (string, error) {
	base := viper.GetString("endpoint")
	slog.Debug("[GenerateURI]", "base", base, "section", section, "command", command)

	if socket, ok := socketPath(base); ok {
		if socket == "" {
			return "", errors.New("endpoint " + UnixSocketScheme + " is missing the socket path")
		}
		base = unixSocketBaseURL
	}

	scheme := ""
	if !strings.Contains(base, "://") {
		scheme = "http://"
	}

	uri := fmt.Sprintf("%s%s/%s/%s",
		scheme,
		base,
		section,
		command,
	)

	var myurl, err = url.Parse(uri)
	if err != nil {
		return "", err
	}

	myurl.Path = path.Clean(myurl.Path)

	res, _ := url.PathUnescape(myurl.String())
	slog.Debug("[GenerateURI] Result", "uri", uri, "url", myurl, "url(string)", res)
	return res, nil
}

// GetJSONRequest returns a request prepared for default JSON responses
func GetJSONRequestTimeout(timeout time.Duration) *resty.Request {
	request := GetRequestTimeout(timeout).
		SetResult(Response{}).
		SetError(Response{})
	if RawJSON {
		request.
			SetDoNotParseResponse(true)
	}
	return request
}

func GetJSONRequest() *resty.Request {
	return GetJSONRequestTimeout(DefaultTimeout)
}

// newClient returns a resty.Client for the Supervisor API. When socket is not
// empty, every connection goes to that Unix socket.
func newClient(timeout time.Duration, socket string) *resty.Client {
	c := resty.New()

	// Default is no timeout. This can lead to lockup the CLI
	// in case the server does not respond. Set a somewhat low
	// timeout for our local only use case.
	c.SetTimeout(timeout)

	if socket != "" {
		// Start from resty's default transport so connection pooling and
		// timeouts match TCP, only the way connections are made changes
		transport := c.GetClient().Transport.(*http.Transport).Clone()
		// The socket is local, never route it through an HTTP proxy
		transport.Proxy = nil
		transport.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, "unix", socket)
		}
		c.SetTransport(transport)
	}

	// Registering Response Middleware
	c.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// explore response object
		slog.Debug("Response", "statuscode", resp.StatusCode(), "status", resp.Status(), "time", resp.Time(), "received-at", resp.ReceivedAt(), "headers", resp.Header(), "request", resp.Request.RawRequest, "body", resp)

		return nil // if its success otherwise return error
	})

	return c
}

// GetRequest returns a resty.Request object prepared for an API call
func GetRequestTimeout(timeout time.Duration) *resty.Request {
	socket, overSocket := socketPath(viper.GetString("endpoint"))

	if client == nil {
		client = newClient(timeout, socket)
	}

	request := client.R().
		SetHeader("Accept", "application/json")

	// Supervisor identifies the caller by the socket, no token needed
	if !overSocket {
		request.SetAuthToken(viper.GetString("api-token"))
	}

	return request
}

// ShowJSONResponse formats a JSON response for human readers
func ShowJSONResponse(resp *resty.Response) (success bool) {
	return ShowJSONResponseTransform(resp, nil)
}

// ShowJSONResponseTransform formats a JSON response for human readers, applying
// transform to the data map of a successful response before rendering it. This
// allows commands to humanize individual fields. In raw JSON mode the transform
// is not applied, as the output is meant for machine consumption.
func ShowJSONResponseTransform(resp *resty.Response, transform func(map[string]any)) (success bool) {
	if RawJSON {
		// when we are returning raw JSON, all handling is for the consumer of the JSON
		success = true
		body := resp.RawBody()
		if b, err := io.ReadAll(body); err == nil {
			fmt.Print(string(b))
		}
		_ = body.Close()
		return
	}
	var data *Response
	if resp.IsSuccess() {
		data = resp.Result().(*Response)
	} else {
		data = resp.Error().(*Response)
	}
	switch data.Result {
	case "ok":
		if len(data.Data) == 0 {
			fmt.Println("Command completed successfully.")
		} else {
			if transform != nil {
				transform(data.Data)
			}
			j, err := json.Marshal(data.Data)
			if err != nil {
				PrintError(err)
				return
			}

			d, err := yaml.JSONToYAML(j)
			if err != nil {
				PrintError(err)
				return
			}
			fmt.Print(string(d))
		}
		success = true
	case "error":
		PrintErrorString(data.Message)
	default:
		d, err := yaml.Marshal(data)
		if err != nil {
			PrintError(err)
			return
		}
		fmt.Print(string(d))
	}
	return
}

func GetRequest() *resty.Request {
	return GetRequestTimeout(DefaultTimeout)
}

// Streams out a text response, like ones from /host/logs
func StreamTextResponse(resp *resty.Response) (success bool) {
	success = true
	for {
		p := make([]byte, 4096)
		n, err := resp.RawBody().Read(p)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			success = false
			break
		}
		fmt.Print(string(p[:n]))
		if err == io.EOF {
			break
		}
	}
	return
}

func AskForConfirmation(prompt string, tries int) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	if tries <= 0 {
		tries = 2
	}

	for ; tries > 0; tries-- {
		fmt.Printf("%s [enter YES to confirm] ", prompt)

		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			return false, err
		}

		// Require YES or yes explicitly to confirm
		res = strings.ToLower(strings.TrimSpace(res))
		if res == "yes" {
			return true, nil
		}

		// If user enters no or n then stop. Else retry since they entered something unknown
		if len(res) > 0 && res[0] == 'n' {
			return false, nil
		}
	}
	return false, nil
}

func ReadInteger(prompt string, tries int, min int, max int) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	if tries <= 0 {
		tries = 2
	}

	for ; tries > 0; tries-- {
		fmt.Printf("%s [%d-%d]: ", prompt, min, max)

		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			return 0, err
		}

		res = strings.TrimSpace(res)
		if len(res) == 0 {
			continue
		}

		val, err := strconv.Atoi(res)
		if err != nil || val < min || val > max {
			fmt.Printf("Invalid value. Must be between %d and %d.\n", min, max)
			continue
		}
		return val, nil
	}

	return 0, errors.New("maximum tries exceeded")
}

func ReadPassword(repeat bool) (string, error) {
	initialState, err := term.GetState(syscall.Stdin)
	if err != nil {
		return "", err
	}

	// Make sure terminal is restored on termination
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		err := term.Restore(syscall.Stdin, initialState)
		if err != nil {
			fmt.Println("Failed to restore terminal state!")
			panic(err)
		}
		os.Exit(1)
	}()

	fmt.Print("Password: ")
	password, err := term.ReadPassword(syscall.Stdin)
	fmt.Println()
	if err != nil {
		return "", err
	}

	if repeat {
		fmt.Print("Password (again): ")
		password2, err := term.ReadPassword(syscall.Stdin)
		fmt.Println()
		if err != nil {
			return "", err
		}

		if string(password) != string(password2) {
			return "", errors.New("passwords do not match")
		}
	}

	return string(password), nil
}

func PrintErrorString(err string) {
	if term.IsTerminal(syscall.Stderr) {
		fmt.Fprintln(os.Stderr, "\033[1;31mError:\033[0m", err)
	} else {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

// PrintError prints an error message with "Error: " prefix for user-friendly output
func PrintError(err error) {
	PrintErrorString(err.Error())
}

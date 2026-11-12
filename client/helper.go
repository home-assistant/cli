package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/home-assistant/cli/config"

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
	base := config.GetString("endpoint")
	slog.Debug("[GenerateURI]", "base", base, "section", section, "command", command)

	if socket, ok := socketPath(base); ok {
		if socket == "" {
			return "", errors.New("endpoint " + UnixSocketScheme + " is missing the socket path")
		}
		base = unixSocketBaseURL
	} else if strings.HasPrefix(base, "unix:/") {
		// A socket path missing a slash, would otherwise be sent over TCP to a
		// host named "unix"
		return "", fmt.Errorf("invalid endpoint %q, use %s/path/to/socket", base, UnixSocketScheme)
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

// GetJSONRequestTimeout returns a request prepared for default JSON responses
func GetJSONRequestTimeout(timeout time.Duration) *Request {
	request := GetRequestTimeout(timeout)
	if RawJSON {
		// The body is passed through as-is, see ShowJSONResponseTransform
		request.stream = true
	} else {
		request.parseJSON = true
	}
	return request
}

func GetJSONRequest() *Request {
	return GetJSONRequestTimeout(DefaultTimeout)
}

// GetRequestTimeout returns a Request object prepared for an API call. The
// timeout applies to this request only, 0 disables it.
func GetRequestTimeout(timeout time.Duration) *Request {
	socket, overSocket := socketPath(config.GetString("endpoint"))

	request := &Request{
		socket:     socket,
		timeout:    timeout,
		header:     http.Header{},
		pathParams: map[string]string{},
	}
	request.SetHeader("Accept", "application/json")
	request.SetHeader("User-Agent", userAgent())

	// Supervisor identifies the caller by the socket, no token needed
	if token := config.GetString("api-token"); !overSocket && token != "" {
		request.SetHeader("Authorization", "Bearer "+token)
	}

	// Lets an app calling the CLI identify itself to Supervisor
	if source := config.GetString("hass-source"); source != "" {
		request.SetHeader("X-Hass-Source", source)
	}

	return request
}

// ShowJSONResponse formats a JSON response for human readers
func ShowJSONResponse(resp *HTTPResponse) (success bool) {
	return ShowJSONResponseTransform(resp, nil)
}

// ShowJSONResponseTransform formats a JSON response for human readers, applying
// transform to the data map of a successful response before rendering it. This
// allows commands to humanize individual fields. In raw JSON mode the transform
// is not applied, as the output is meant for machine consumption. Both modes
// decide success the same way, from the "result" field of the response body.
func ShowJSONResponseTransform(resp *HTTPResponse, transform func(map[string]any)) (success bool) {
	if RawJSON {
		// The client does not parse the response in raw JSON mode. The body is
		// passed through to the consumer as-is, only its "result" field is
		// looked at to decide the outcome.
		body := resp.RawBody()
		b, err := io.ReadAll(body)
		_ = body.Close()
		if err != nil {
			PrintError(err)
			return
		}
		fmt.Print(string(b))
		var result struct {
			Result string `json:"result"`
		}
		if err := json.Unmarshal(b, &result); err != nil {
			PrintError(err)
			return
		}
		return result.Result == "ok"
	}
	var data *Response
	if resp.IsSuccess() {
		data = resp.Result()
	} else {
		data = resp.Error()
	}
	switch data.Result {
	case "ok":
		if len(data.Data) == 0 {
			fmt.Println("Command completed successfully.")
		} else {
			if transform != nil {
				transform(data.Data)
			}
			d, err := toYAML(data.Data)
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
		d, err := toYAML(data)
		if err != nil {
			PrintError(err)
			return
		}
		fmt.Print(string(d))
	}
	return
}

func GetRequest() *Request {
	return GetRequestTimeout(DefaultTimeout)
}

// Streams out a text response, like ones from /host/logs
func StreamTextResponse(resp *HTTPResponse) (success bool) {
	body := resp.RawBody()
	defer body.Close()
	if _, err := io.Copy(os.Stdout, body); err != nil {
		fmt.Println(err)
		return false
	}
	return true
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

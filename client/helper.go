package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/term"
	"io"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	yaml "github.com/ghodss/yaml"
	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"strings"
)

const DefaultTimeout = 30 * time.Second
const ContainerOperationTimeout = 10 * time.Minute
const ContainerDownloadTimeout = 1 * time.Hour
const OsDownloadTimeout = 1 * time.Hour
const BackupTimeout = 3 * time.Hour
const RebootTimeout = 90 * time.Second

var client *resty.Client

// Response is the default JSON response from the Home Assistant Supervisor
type Response struct {
	Result  string         `json:"result"`
	Message string         `json:"message,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

// URLHelper returns a URL built from the arguments
func URLHelper(section, command string) (string, error) {
	base := viper.GetString("endpoint")
	log.WithFields(log.Fields{
		"base":    base,
		"section": section,
		"command": command,
	}).Debug("[GenerateURI]")

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
	log.WithFields(log.Fields{
		"uri":         uri,
		"url":         myurl,
		"url(string)": res,
	}).Debug("[GenerateURI] Result")
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

// GetRequest returns a resty.Request object prepared for an API call
func GetRequestTimeout(timeout time.Duration) *resty.Request {
	apiToken := viper.GetString("api-token")

	if client == nil {
		client = resty.New()

		// Default is no timeout. This can lead to lockup the CLI
		// in case the server does not respond. Set a somewhat low
		// timeout for our local only use case.
		client.SetTimeout(timeout)

		// Registering Response Middleware
		client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			// explore response object
			log.WithFields(log.Fields{
				"statuscode":  resp.StatusCode(),
				"status":      resp.Status(),
				"time":        resp.Time(),
				"received-at": resp.ReceivedAt(),
				"headers":     resp.Header(),
				"request":     resp.Request.RawRequest,
				"body":        resp,
			}).Debug("Response")

			return nil // if its success otherwise return error
		})
	}

	return client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiToken)
}

// ShowJSONResponse formats a JSON response for human readers
func ShowJSONResponse(resp *resty.Response) (success bool) {
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
		success = true
		if len(data.Data) == 0 {
			fmt.Println("Command completed successfully.")
		} else {
			j, err := json.Marshal(data.Data)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			d, err := yaml.JSONToYAML(j)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			fmt.Print(string(d))
		}
	case "error":
		fmt.Printf("Error: %s\n", data.Message)
	default:
		d, err := yaml.Marshal(data)
		if err != nil {
			log.Fatalf("error: %v", err)
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

func AskForConfirmation(prompt string, tries int) bool {
	reader := bufio.NewReader(os.Stdin)
	if tries <= 0 {
		tries = 2
	}

	for ; tries > 0; tries-- {
		fmt.Printf("%s [enter YES to confirm] ", prompt)

		res, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error: %v", err)
			continue
		}

		// Require YES or yes explicitly to confirm
		res = strings.ToLower(strings.TrimSpace(res))
		if res == "yes" {
			return true
		}

		// If user enters no or n then stop. Else retry since they entered something unknown
		if len(res) > 0 && res[0] == 'n' {
			return false
		}
	}
	return false
}

func ReadInteger(prompt string, tries int, min int, max int) (bool, int) {
	reader := bufio.NewReader(os.Stdin)
	if tries <= 0 {
		tries = 2
	}

	for ; tries > 0; tries-- {
		fmt.Printf("%s [%d-%d]: ", prompt, min, max)

		res, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error: %v", err)
			continue
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
		return true, val
	}

	return false, -1
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
		password2, err := term.ReadPassword(int(os.Stdin.Fd()))
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

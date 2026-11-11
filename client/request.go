package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Request is a Supervisor API request. It is created with GetRequest or
// GetJSONRequest (or their Timeout variants), configured with the Set methods
// and sent with Get, Post, Put, Delete or Send.
type Request struct {
	// URL is the URL used by Send. It may contain {name} placeholders, see
	// SetPathParam.
	URL string

	client     *http.Client
	timeout    time.Duration
	header     http.Header
	pathParams map[string]string
	body       any
	parseJSON  bool
	stream     bool
}

// SetHeader sets a request header, replacing any previous value.
func (r *Request) SetHeader(key, value string) *Request {
	r.header.Set(key, value)
	return r
}

// SetPathParam sets the value replacing the {name} placeholder in the URL.
// The value is escaped, so it always ends up as a single path segment.
func (r *Request) SetPathParam(name, value string) *Request {
	r.pathParams[name] = value
	return r
}

// SetPathParams sets several path parameters, see SetPathParam.
func (r *Request) SetPathParams(params map[string]string) *Request {
	for name, value := range params {
		r.SetPathParam(name, value)
	}
	return r
}

// SetBody sets the request body, which is sent JSON encoded.
func (r *Request) SetBody(body any) *Request {
	r.body = body
	return r
}

// SetDoNotParseResponse makes the request leave the response body unread.
// The caller reads it from HTTPResponse.RawBody and has to close it. Use this
// for responses which are streamed, like logs.
func (r *Request) SetDoNotParseResponse(stream bool) *Request {
	r.stream = stream
	return r
}

// Get sends the request with the GET method.
func (r *Request) Get(url string) (*HTTPResponse, error) {
	return r.execute(http.MethodGet, url)
}

// Post sends the request with the POST method.
func (r *Request) Post(url string) (*HTTPResponse, error) {
	return r.execute(http.MethodPost, url)
}

// Put sends the request with the PUT method.
func (r *Request) Put(url string) (*HTTPResponse, error) {
	return r.execute(http.MethodPut, url)
}

// Delete sends the request with the DELETE method.
func (r *Request) Delete(url string) (*HTTPResponse, error) {
	return r.execute(http.MethodDelete, url)
}

// Send sends the request with the GET method to the URL field.
func (r *Request) Send() (*HTTPResponse, error) {
	return r.execute(http.MethodGet, r.URL)
}

func (r *Request) execute(method, rawURL string) (*HTTPResponse, error) {
	for name, value := range r.pathParams {
		rawURL = strings.ReplaceAll(rawURL, "{"+name+"}", url.PathEscape(value))
	}

	var body io.Reader
	if r.body != nil {
		b, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}

	// The timeout covers the whole exchange, including reading the body. It
	// is set per request, so each request gets the timeout it asked for.
	ctx, cancel := context.Background(), context.CancelFunc(func() {})
	if r.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		cancel()
		return nil, err
	}
	req.Header = r.header.Clone()
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	start := time.Now()
	raw, err := r.client.Do(req)
	if err != nil {
		cancel()
		return nil, err
	}

	resp := &HTTPResponse{raw: raw, result: &Response{}, errResult: &Response{}}

	if r.stream {
		// The caller reads the body, keep the context alive until it is closed
		raw.Body = &cancelOnClose{ReadCloser: raw.Body, cancel: cancel}
		slog.Debug("Response", "method", method, "url", rawURL, "status", raw.Status,
			"time", time.Since(start), "headers", raw.Header)
		return resp, nil
	}
	defer cancel()

	resp.body, err = io.ReadAll(raw.Body)
	_ = raw.Body.Close()
	if err != nil {
		return nil, err
	}
	raw.Body = io.NopCloser(bytes.NewReader(resp.body))

	slog.Debug("Response", "method", method, "url", rawURL, "status", raw.Status,
		"time", time.Since(start), "headers", raw.Header, "body", string(resp.body))

	if r.parseJSON && len(resp.body) > 0 && IsJSONType(raw.Header.Get("Content-Type")) {
		switch {
		case resp.IsSuccess():
			if err := json.Unmarshal(resp.body, resp.result); err != nil {
				return nil, err
			}
		case resp.IsError():
			if err := json.Unmarshal(resp.body, resp.errResult); err != nil {
				slog.Debug("Cannot parse error response", "error", err)
			}
		}
	}

	return resp, nil
}

// cancelOnClose releases the request context once the body is closed.
type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c *cancelOnClose) Close() error {
	err := c.ReadCloser.Close()
	c.cancel()
	return err
}

// HTTPResponse is the response to a Request.
type HTTPResponse struct {
	raw       *http.Response
	body      []byte
	result    *Response
	errResult *Response
}

// StatusCode returns the HTTP status code, e.g. 200.
func (r *HTTPResponse) StatusCode() int {
	return r.raw.StatusCode
}

// Status returns the HTTP status line, e.g. "200 OK".
func (r *HTTPResponse) Status() string {
	return r.raw.Status
}

// Header returns the response headers.
func (r *HTTPResponse) Header() http.Header {
	return r.raw.Header
}

// IsSuccess reports if the status code is in the 2xx range.
func (r *HTTPResponse) IsSuccess() bool {
	return r.raw.StatusCode >= 200 && r.raw.StatusCode <= 299
}

// IsError reports if the status code is 400 or above.
func (r *HTTPResponse) IsError() bool {
	return r.raw.StatusCode >= 400
}

// Body returns the response body. It is empty when the request was sent with
// SetDoNotParseResponse, use RawBody then.
func (r *HTTPResponse) Body() []byte {
	return r.body
}

// RawBody returns the response body as a stream. When the request was sent
// with SetDoNotParseResponse the caller has to close it.
func (r *HTTPResponse) RawBody() io.ReadCloser {
	return r.raw.Body
}

// Result returns the parsed JSON body of a successful response to a JSON
// request. It is never nil, it is empty when there was nothing to parse.
func (r *HTTPResponse) Result() *Response {
	return r.result
}

// Error returns the parsed JSON body of an error response to a JSON request.
// It is never nil, it is empty when there was nothing to parse.
func (r *HTTPResponse) Error() *Response {
	return r.errResult
}

// IsJSONType reports if a Content-Type header value denotes JSON, e.g.
// application/json or application/problem+json.
func IsJSONType(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	kind, subtype, _ := strings.Cut(mediaType, "/")
	return (kind == "application" || kind == "text") && strings.Contains(subtype, "json")
}

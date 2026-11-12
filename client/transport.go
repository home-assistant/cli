package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// roundTrip sends req over a connection of its own and returns the response.
// When socket is not empty the connection goes to that Unix socket, else to
// the host of the request URL. The deadline covers the whole exchange,
// including reading the response body, the zero value disables it. Closing
// the response body closes the connection.
//
// This is used instead of http.Client, as http.Transport links TLS and HTTP/2
// support into the binary even if it is never used. The Supervisor API is
// plain HTTP/1.1, so the request is written to the connection directly.
// Unlike http.Transport this does not reuse connections, honor proxy settings
// or follow redirects.
func roundTrip(req *http.Request, socket string, deadline time.Time) (*http.Response, error) {
	if req.URL.Scheme != "http" {
		return nil, fmt.Errorf("unsupported protocol scheme %q, the endpoint has to be plain HTTP or a Unix socket", req.URL.Scheme)
	}

	network, address := "tcp", req.URL.Host
	if socket != "" {
		network, address = "unix", socket
	} else if req.URL.Port() == "" {
		address = net.JoinHostPort(req.URL.Hostname(), "80")
	}

	dialer := net.Dialer{Deadline: deadline}
	conn, err := dialer.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if err := conn.SetDeadline(deadline); err != nil {
		_ = conn.Close()
		return nil, err
	}

	// The connection is not reused, let the server know
	req.Close = true
	if err := req.Write(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	resp.Body = &connBody{Reader: resp.Body, conn: conn}
	return resp, nil
}

// connBody is a response body which closes its connection. The body itself
// is not closed, as that reads it to the end first, which never comes for a
// response streamed forever like followed logs.
type connBody struct {
	io.Reader
	conn net.Conn
}

func (b *connBody) Close() error {
	return b.conn.Close()
}

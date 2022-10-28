package agent

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// customHttpError for check results if request contains status 4** or 5**
type customHTTPError struct {
	Message string
	Status  int
}

func (e *customHTTPError) Error() string {
	return fmt.Sprintf("request ended with status %d and error: %s", e.Status, e.Message)
}

// Client http implementation
type client struct {
	http.Client
}

// NewClient constructor for client
func NewClient(timeout time.Duration, maxIdleConns int) Client {
	transport := &http.Transport{}

	if maxIdleConns != 0 {
		transport.MaxIdleConns = maxIdleConns
	}

	newClient := &client{}
	newClient.Transport = transport

	if timeout != 0 {
		newClient.Timeout = timeout
	}

	return newClient

}

// DoRequest sending request to resource
func (c *client) DoRequest(method, url string, headers map[string]string, body []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		var b []byte

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return &customHTTPError{
			Message: string(b),
			Status:  resp.StatusCode,
		}
	}

	return nil
}

// Shutdown gracefully shutdown for correct close cached connections
func (c *client) Shutdown() {
	c.CloseIdleConnections()
}

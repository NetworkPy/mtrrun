package agent

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrNotStatusOK = fmt.Errorf("response status not 200")
)

// Client http implementation
type client struct {
	http.Client
}

// NewClient constructor for client
func NewClient(timeout int, maxIdleConns int) Client {
	transport := &http.Transport{}

	if maxIdleConns != 0 {
		transport.MaxIdleConns = maxIdleConns
	}

	newClient := &client{}
	newClient.Transport = transport

	if timeout != 0 {
		newClient.Timeout = time.Duration(timeout) * time.Second
	}

	return newClient

}

// DoRequest sending request to resource
func (c *client) DoRequest(method, url string, header map[string]string, body []byte) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	for k, v := range header {
		req.Header.Add(k, v)
	}

	if err != nil {
		return err
	}

	resp, err := c.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrNotStatusOK
	}

	return nil
}

// Shutdown gracefully shutdown for correct close cached connections
func (c *client) Shutdown() {
	c.CloseIdleConnections()
}

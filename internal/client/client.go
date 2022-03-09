package client

import (
	"time"

	"github.com/gojektech/heimdall/httpclient"
)

const (
	defaultHTTPClientTimeout = 5 * time.Second
	defaultHTTPClientRetries = 3
)

// Client is a wrapper for the heimdall client
type Client struct {
	*httpclient.Client
}

// NewDefaultClient returns httpclient instance with default config
func NewDefaultClient() *Client {
	return &Client{
		Client: httpclient.NewClient(
			httpclient.WithHTTPTimeout(defaultHTTPClientTimeout),
			httpclient.WithRetryCount(defaultHTTPClientRetries),
		),
	}
}

// NewCustomClient returns httpclient instance with given custom config
func NewCustomClient(retries int, timeout time.Duration) *Client {
	return &Client{
		Client: httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetryCount(retries),
		),
	}
}

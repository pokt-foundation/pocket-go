package client

import (
	"time"

	"github.com/gojektech/heimdall"
	"github.com/gojektech/heimdall/httpclient"
)

const (
	defaultHTTPClientTimeout = 5 * time.Second
	defaultHTTPClientRetries = 3

	initialBackoffTimeout = 2 * time.Millisecond
	maxBackoffTimeout     = 9 * time.Millisecond
	exponentFactor        = 2
	maxJitterInterval     = 2 * time.Millisecond
)

var (
	backoff = heimdall.NewExponentialBackoff(initialBackoffTimeout, maxBackoffTimeout, exponentFactor, maxJitterInterval)
	retrier = heimdall.NewRetrier(backoff)
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
			httpclient.WithRetrier(retrier),
		),
	}
}

// NewCustomClient returns httpclient instance with given custom config
func NewCustomClient(retries int, timeout time.Duration) *Client {
	return &Client{
		Client: httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetryCount(retries),
			httpclient.WithRetrier(retrier),
		),
	}
}

package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

// PostWithURLJSONParams does post request with JSON param
func (client *Client) PostWithURLJSONParams(url string, params interface{}, headers http.Header) (*http.Response, error) {
	var body io.Reader

	if params != nil {
		rawBody, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBufferString(string(rawBody))
	}

	// Needed header for JSON request
	headers.Set("Content-Type", "application/json")

	return client.Post(url, body, headers)
}

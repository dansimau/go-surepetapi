package httpclient

import "github.com/dansimau/go-surepetapi/pkg/httpdebug"

type Option func(c *Client)

func WithDebugLogging(enabled bool) Option {
	return func(c *Client) {
		if enabled {
			httpdebug.EnableDebugLogging(c.retryClient.HTTPClient)
		}
	}
}

func WithDefaultHeaders(headers map[string]string) Option {
	return func(c *Client) {
		c.defaultHeaders = headers
	}
}

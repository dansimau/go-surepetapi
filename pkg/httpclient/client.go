package httpclient

import (
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

var _ Interface = (*Client)(nil)

type Client struct {
	defaultHeaders map[string]string
	retryClient    *retryablehttp.Client

	BaseURL string
}

type Interface interface {
	Request(r RequestParams) (*http.Response, error)
	SetFollowRedirects(followRedirects bool)
	StandardHTTPClient() *http.Client
}

func New(opts ...Option) *Client {
	retryClient := retryablehttp.NewClient()

	retryClient.Logger = nil

	// these are hardcoded now but can be put into config/options later
	retryClient.RetryMax = 2
	retryClient.HTTPClient.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Second * 3,
		}).Dial,
		TLSHandshakeTimeout: time.Second * 3,
	}
	retryClient.HTTPClient.Timeout = time.Second * 15

	c := &Client{
		retryClient: retryClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetFollowRedirects sets whether or not this client should follow redirects.
func (c *Client) SetFollowRedirects(followRedirects bool) {
	c.retryClient.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}

// StandardHTTPClient returns a HTTPClient that can be used for compatibility
// with other libraries.
func (c *Client) StandardHTTPClient() *http.Client {
	return c.retryClient.StandardClient()
}

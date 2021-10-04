package httpclient

import (
	"encoding/json"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type RequestParams struct {
	BodyJSON interface{}
	Headers  map[string]string
	Method   string
	Params   map[string]string
	URL      string
}

// Request makes a request to the test app and returns the response.
func (c *Client) Request(r RequestParams) (*http.Response, error) {
	if c.BaseURL != "" {
		r.URL = c.BaseURL + r.URL
	}

	var body interface{}
	if r.BodyJSON != nil {
		b, err := json.Marshal(r.BodyJSON)
		if err != nil {
			return nil, err
		}

		body = b
	}

	req, err := retryablehttp.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}

	for key, val := range c.defaultHeaders {
		req.Header.Add(key, val)
	}

	for key, val := range r.Headers {
		req.Header.Add(key, val)
	}

	if r.BodyJSON != nil {
		req.Header.Add("content-type", "application/json")
	}

	if r.Params != nil {
		query := req.URL.Query()
		for key, val := range r.Params {
			query.Add(key, val)
		}
		req.URL.RawQuery = query.Encode()
	}

	return c.retryClient.Do(req)
}

package httpclient

import "net/http"

var DefaultClient *Client

func init() {
	DefaultClient = New()
}

func Request(r RequestParams) (*http.Response, error) {
	return DefaultClient.Request(r)
}

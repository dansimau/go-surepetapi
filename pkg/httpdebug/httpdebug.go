package httpdebug

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

type loggingRoundTripper struct {
	proxied http.RoundTripper
}

func (lrt *loggingRoundTripper) dumpResponse(resp *http.Response) {
	println("\n--- DEBUG: HTTP response:")
	defer println("\n---")

	bytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Unable to dump response header: %v\n", err)
		return
	}

	fmt.Fprint(os.Stderr, string(bytes))
}

func (lrt *loggingRoundTripper) dumpRequest(req *http.Request) {
	println("\n--- DEBUG: HTTP request:")
	defer println("\n---")

	bytes, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Unable to dump request: %v\n", err)
		return
	}

	fmt.Fprint(os.Stderr, string(bytes))
}

func (lrt loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	lrt.dumpRequest(req)

	// Send the request, get the response (or the error)
	res, err := lrt.proxied.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	lrt.dumpResponse(res)

	return res, err
}

// EnableDebugLogging turns on logging of requests and responses to stderr for
// the specified HTTP client.
func EnableDebugLogging(client *http.Client) {
	client.Transport = loggingRoundTripper{client.Transport}
}

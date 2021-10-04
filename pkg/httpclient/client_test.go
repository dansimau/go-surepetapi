package httpclient_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dansimau/go-surepetapi/pkg/httpclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	testContent := []byte(`{"hello": "world"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(testContent)
	}))
	defer server.Close()

	client := httpclient.New()

	res, err := client.Request(httpclient.RequestParams{
		Method: "GET",
		URL:    server.URL + "/",
	})
	require.NoError(t, err)

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, testContent, b)
}

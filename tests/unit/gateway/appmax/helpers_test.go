package appmax_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	appmax "github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/gateway/appmax"
	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func clientWithTransport(t *testing.T, fn roundTripperFunc) *appmax.Client {
	t.Helper()

	httpClient := &http.Client{Transport: fn}
	client, err := appmax.NewClientWithHTTPClient("https://auth.example.com", "https://api.example.com", httpClient)
	require.NoError(t, err)
	return client
}

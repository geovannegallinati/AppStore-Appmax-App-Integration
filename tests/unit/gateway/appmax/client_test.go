package appmax_test

import (
	"net/http"
	"testing"
	"time"

	appmax "github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/gateway/appmax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientWithOptions_UsesDefaults(t *testing.T) {
	client, err := appmax.NewClientWithOptions("https://auth.example.com", "https://api.example.com", appmax.ClientOptions{})

	require.NoError(t, err)
	assert.Equal(t, appmax.DefaultHTTPTimeout, client.HTTPClientTimeout())
	assert.Equal(t, 0, client.RetryMax())
	assert.Equal(t, appmax.DefaultRetryWait, client.RetryWait())
}

func TestNewClientWithOptions_UsesCustomValues(t *testing.T) {
	transport := roundTripperFunc(func(_ *http.Request) (*http.Response, error) {
		return nil, nil
	})

	client, err := appmax.NewClientWithOptions("https://auth.example.com", "https://api.example.com", appmax.ClientOptions{
		Timeout:   5 * time.Second,
		RetryMax:  2,
		RetryWait: 50 * time.Millisecond,
		Transport: transport,
	})

	require.NoError(t, err)
	assert.Equal(t, 5*time.Second, client.HTTPClientTimeout())
	assert.Equal(t, 2, client.RetryMax())
	assert.Equal(t, 50*time.Millisecond, client.RetryWait())
	assert.NotNil(t, client.HTTPClientTransport())
}

func TestNewClientWithOptions_RejectsInvalidURLs(t *testing.T) {
	client, err := appmax.NewClient("", "")

	require.Error(t, err)
	assert.Nil(t, client)
}

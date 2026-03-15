package appmax

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewClientWithOptions_UsesDefaults(t *testing.T) {
	client, err := NewClientWithOptions("https://auth.example.com", "https://api.example.com", ClientOptions{})

	require.NoError(t, err)
	assert.Equal(t, defaultHTTPTimeout, client.httpClient.Timeout)
	assert.Equal(t, 0, client.retryMax)
	assert.Equal(t, defaultRetryWait, client.retryWait)
}

func TestNewClientWithOptions_UsesCustomValues(t *testing.T) {
	transport := roundTripperFunc(func(_ *http.Request) (*http.Response, error) {
		return nil, nil
	})

	client, err := NewClientWithOptions("https://auth.example.com", "https://api.example.com", ClientOptions{
		Timeout:   5 * time.Second,
		RetryMax:  2,
		RetryWait: 50 * time.Millisecond,
		Transport: transport,
	})

	require.NoError(t, err)
	assert.Equal(t, 5*time.Second, client.httpClient.Timeout)
	assert.Equal(t, 2, client.retryMax)
	assert.Equal(t, 50*time.Millisecond, client.retryWait)
	assert.NotNil(t, client.httpClient.Transport)
}

func TestNewClientWithOptions_RejectsInvalidURLs(t *testing.T) {
	client, err := NewClient("", "")

	require.Error(t, err)
	assert.Nil(t, client)
}

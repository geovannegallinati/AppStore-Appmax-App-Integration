package appmax

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	gatewaycontracts "github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/gateway/appmax/contracts"
)

const (
	DefaultHTTPTimeout = 90 * time.Second
	DefaultRetryWait   = 200 * time.Millisecond
)

type ClientOptions struct {
	HTTPClient    *http.Client
	Timeout       time.Duration
	Transport     http.RoundTripper
	RetryMax      int
	RetryWait     time.Duration
	RetryStatuses []int
}

type Client struct {
	httpClient    *http.Client
	authBaseURL   string
	apiBaseURL    string
	retryMax      int
	retryWait     time.Duration
	retryStatuses map[int]struct{}
}

var _ gatewaycontracts.Gateway = (*Client)(nil)

func NewClient(authURL, apiURL string) (*Client, error) {
	return NewClientWithOptions(authURL, apiURL, ClientOptions{})
}

func NewClientWithHTTPClient(authURL, apiURL string, httpClient *http.Client) (*Client, error) {
	return NewClientWithOptions(authURL, apiURL, ClientOptions{HTTPClient: httpClient})
}

func NewClientWithOptions(authURL, apiURL string, options ClientOptions) (*Client, error) {
	if strings.TrimSpace(authURL) == "" {
		return nil, fmt.Errorf("new appmax client: auth url is required")
	}
	if strings.TrimSpace(apiURL) == "" {
		return nil, fmt.Errorf("new appmax client: api url is required")
	}

	httpClient := buildHTTPClient(options)
	retryWait := options.RetryWait
	if retryWait <= 0 {
		retryWait = DefaultRetryWait
	}
	retryMax := options.RetryMax
	if retryMax < 0 {
		retryMax = 0
	}

	retryStatuses := make(map[int]struct{}, len(options.RetryStatuses))
	for _, s := range options.RetryStatuses {
		retryStatuses[s] = struct{}{}
	}

	return &Client{
		httpClient:    httpClient,
		authBaseURL:   strings.TrimRight(authURL, "/"),
		apiBaseURL:    strings.TrimRight(apiURL, "/"),
		retryMax:      retryMax,
		retryWait:     retryWait,
		retryStatuses: retryStatuses,
	}, nil
}

func buildHTTPClient(options ClientOptions) *http.Client {
	client := options.HTTPClient
	if client == nil {
		client = &http.Client{}
	}

	timeout := options.Timeout
	if timeout <= 0 {
		timeout = DefaultHTTPTimeout
	}

	mutated := false
	if client.Timeout == 0 || options.Timeout > 0 {
		client = cloneHTTPClient(client)
		client.Timeout = timeout
		mutated = true
	}

	if options.Transport != nil {
		if !mutated {
			client = cloneHTTPClient(client)
			mutated = true
		}
		client.Transport = options.Transport
	}

	return client
}

func cloneHTTPClient(client *http.Client) *http.Client {
	clone := *client
	return &clone
}

func (c *Client) RetryMax() int {
	return c.retryMax
}

func (c *Client) RetryWait() time.Duration {
	return c.retryWait
}

func (c *Client) HTTPClientTimeout() time.Duration {
	return c.httpClient.Timeout
}

func (c *Client) HTTPClientTransport() http.RoundTripper {
	return c.httpClient.Transport
}

package appmax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
)

type upstreamHTTPError struct {
	statusCode   int
	traceHeaders string
	body         string
	message      string
}

func (e *upstreamHTTPError) Error() string {
	return fmt.Sprintf("unexpected status %d%s: %s", e.statusCode, e.traceHeaders, e.body)
}

func (e *upstreamHTTPError) HTTPStatus() int {
	return e.statusCode
}

func (e *upstreamHTTPError) UpstreamMessage() string {
	return e.message
}

func (c *Client) do(ctx context.Context, method, endpoint string, body any, bearerToken string) (*http.Response, error) {
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
	}

	attempts := c.retryMax + 1
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		var reqBody io.Reader
		if payload != nil {
			reqBody = bytes.NewReader(payload)
		}

		req, err := http.NewRequestWithContext(ctx, method, endpoint, reqBody)
		if err != nil {
			return nil, fmt.Errorf("new request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if bearerToken != "" {
			req.Header.Set("Authorization", "Bearer "+bearerToken)
		}

		facades.Log().Debugf("[appmax] → %s %s", method, endpoint)

		resp, err := c.httpClient.Do(req)
		if err == nil {
			facades.Log().Debugf("[appmax] ← %s %s status=%d", method, endpoint, resp.StatusCode)
			if _, retryable := c.retryStatuses[resp.StatusCode]; retryable && attempt < attempts {
				resp.Body.Close()
				facades.Log().Debugf("[appmax] retrying after status %d (attempt %d/%d)", resp.StatusCode, attempt, attempts)
				if !waitWithContext(ctx, c.retryWait) {
					return nil, fmt.Errorf("do request: %w", ctx.Err())
				}
				continue
			}
			return resp, nil
		}

		lastErr = err
		if attempt < attempts {
			if !waitWithContext(ctx, c.retryWait) {
				return nil, fmt.Errorf("do request: %w", ctx.Err())
			}
		}
	}

	return nil, fmt.Errorf("do request: %w", lastErr)
}

func doAndDecode[T any](c *Client, ctx context.Context, method, endpoint string, body any, bearerToken string, expected ...int) (T, error) {
	var out T

	resp, err := c.do(ctx, method, endpoint, body, bearerToken)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp, expected...); err != nil {
		return out, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return out, fmt.Errorf("decode response: %w", err)
	}

	return out, nil
}

func checkStatus(resp *http.Response, expected ...int) error {
	for _, code := range expected {
		if resp.StatusCode == code {
			return nil
		}
	}

	data, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(data))

	traceHeaders := ""
	for _, h := range []string{"CF-Ray", "X-Request-Id", "X-Trace-Id"} {
		if v := resp.Header.Get(h); v != "" {
			traceHeaders += fmt.Sprintf(" %s=%s", h, v)
		}
	}

	body := strings.TrimSpace(string(data))
	message := parseUpstreamMessageFromBody(data)
	err := &upstreamHTTPError{
		statusCode:   resp.StatusCode,
		traceHeaders: traceHeaders,
		body:         body,
		message:      message,
	}
	facades.Log().Errorf("[appmax] %v", err)
	return err
}

func parseUpstreamMessageFromBody(data []byte) string {
	var payload struct {
		Message string `json:"message"`
		Errors  struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return ""
	}

	if payload.Errors.Message != "" {
		return payload.Errors.Message
	}

	return payload.Message
}

func waitWithContext(ctx context.Context, wait time.Duration) bool {
	if wait <= 0 {
		return true
	}

	timer := time.NewTimer(wait)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

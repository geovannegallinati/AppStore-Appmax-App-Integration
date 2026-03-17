package controllers

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const RefundFailedMessage = "refund request failed"

var unexpectedStatusPattern = regexp.MustCompile(`unexpected status (\d{3})`)

type upstreamStatusError interface {
	HTTPStatus() int
}

type upstreamMessageError interface {
	UpstreamMessage() string
}

func UpstreamErrorStatus(err error, fallback int) int {
	if err == nil {
		return fallback
	}

	var statusErr upstreamStatusError
	if errors.As(err, &statusErr) {
		if code := statusErr.HTTPStatus(); code >= 100 && code <= 599 {
			return code
		}
	}

	if code := extractStatusCode(err.Error()); code >= 100 && code <= 599 {
		return code
	}

	return fallback
}

func UpstreamErrorMessage(err error, fallback string) string {
	if err == nil {
		return fallback
	}

	var messageErr upstreamMessageError
	if errors.As(err, &messageErr) {
		if message := strings.TrimSpace(messageErr.UpstreamMessage()); message != "" {
			return message
		}
	}

	message := extractUpstreamMessage(err.Error())
	if message == "" {
		return fallback
	}

	return message
}

func RefundErrorMessage(err error) string {
	return UpstreamErrorMessage(err, RefundFailedMessage)
}

func extractUpstreamMessage(raw string) string {
	start := strings.Index(raw, "{")
	for start >= 0 && start < len(raw) {
		candidate := strings.TrimSpace(raw[start:])
		if message := parseUpstreamMessageJSON(candidate); message != "" {
			return message
		}

		next := strings.Index(raw[start+1:], "{")
		if next < 0 {
			break
		}
		start = start + 1 + next
	}

	return ""
}

func parseUpstreamMessageJSON(payload string) string {
	var objectPayload struct {
		Message string `json:"message"`
		Errors  struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal([]byte(payload), &objectPayload); err != nil {
		return ""
	}

	if objectPayload.Errors.Message != "" {
		return objectPayload.Errors.Message
	}

	return objectPayload.Message
}

func extractStatusCode(raw string) int {
	matches := unexpectedStatusPattern.FindStringSubmatch(raw)
	if len(matches) != 2 {
		return 0
	}

	code, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	return code
}

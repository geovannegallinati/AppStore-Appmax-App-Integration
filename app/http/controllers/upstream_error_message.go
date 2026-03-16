package controllers

import (
	"encoding/json"
	"strings"
)

const refundFailedMessage = "refund request failed"

func upstreamErrorMessage(err error, fallback string) string {
	if err == nil {
		return fallback
	}

	message := extractUpstreamMessage(err.Error())
	if message == "" {
		return fallback
	}

	return message
}

func refundErrorMessage(err error) string {
	return upstreamErrorMessage(err, refundFailedMessage)
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

package requests

import "encoding/json"

type WebhookEnvelopeRequest struct {
	Event     string          `json:"event"`
	EventType string          `json:"event_type"`
	Data      json.RawMessage `json:"data"`
}

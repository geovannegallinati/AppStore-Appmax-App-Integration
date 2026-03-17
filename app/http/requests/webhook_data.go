package requests

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type OptionalInt struct {
	value *int
}

func (o *OptionalInt) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		o.value = nil
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err == nil {
		parsed, err := number.Int64()
		if err == nil {
			value := int(parsed)
			o.value = &value
			return nil
		}
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		str = strings.TrimSpace(str)
		if str == "" {
			o.value = nil
			return nil
		}

		value, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		o.value = &value
		return nil
	}

	return fmt.Errorf("invalid integer value: %s", raw)
}

func (o OptionalInt) Ptr() *int {
	return o.value
}

type WebhookDataRequest struct {
	ID                 OptionalInt      `json:"id"`
	CustomerID         OptionalInt      `json:"customer_id"`
	OrderID            OptionalInt      `json:"order_id"`
	OrderTotalProducts *float64         `json:"order_total_products"`
	Meta               *json.RawMessage `json:"meta"`
}

// ExtractOrderID returns the order ID from the payload using model-aware logic:
//   - Two-Level Flat / Custom Content / Old Legacy: uses data.order_id
//   - Standard / Standard with Meta: uses data.id when data.customer_id is also
//     present, confirming this is an order event (not a customer event where data.id
//     is the customer ID)
//   - Customer events and subscription client-structure events: returns nil
func (d WebhookDataRequest) ExtractOrderID() *int {
	if d.OrderID.Ptr() != nil {
		return d.OrderID.Ptr()
	}
	if d.ID.Ptr() != nil && d.CustomerID.Ptr() != nil {
		return d.ID.Ptr()
	}
	return nil
}

// DetectModel identifies the Appmax webhook payload model from the parsed data
// and the envelope event_type field.
func (d WebhookDataRequest) DetectModel(eventType string) string {
	if eventType == "order" {
		return "old_legacy"
	}
	if d.ID.Ptr() != nil && d.CustomerID.Ptr() != nil {
		if d.Meta != nil {
			return "standard_meta"
		}
		return "standard"
	}
	if d.OrderID.Ptr() != nil {
		if d.OrderTotalProducts != nil {
			return "two_level_flat"
		}
		return "custom_content"
	}
	return "standard"
}

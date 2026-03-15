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
	OrderID OptionalInt `json:"order_id"`
}

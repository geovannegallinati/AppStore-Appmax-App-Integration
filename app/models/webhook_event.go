package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type JSONMap map[string]any

func (j JSONMap) Value() (driver.Value, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (j *JSONMap) Scan(value any) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("JSONMap.Scan: unsupported type %T", value)
	}
	return json.Unmarshal(bytes, j)
}

type WebhookEvent struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Event         string     `gorm:"column:event;not null"`
	EventType     string     `gorm:"column:event_type;not null"`
	AppmaxOrderID *int       `gorm:"column:appmax_order_id"`
	Payload       JSONMap    `gorm:"column:payload;type:jsonb;not null;serializer:json"`
	Processed     bool       `gorm:"column:processed;not null;default:false"`
	ProcessedAt   *time.Time `gorm:"column:processed_at"`
	ErrorMessage  string     `gorm:"column:error_message"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null"`
}

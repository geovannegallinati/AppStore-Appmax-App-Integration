package models

import "time"

type Installation struct {
	ID                   int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ExternalKey          string    `gorm:"column:external_key;not null;uniqueIndex"`
	AppID                string    `gorm:"column:app_id;not null"`
	MerchantClientID     string    `gorm:"column:merchant_client_id;not null"`
	MerchantClientSecret string    `gorm:"column:merchant_client_secret;not null"`
	ExternalID           string    `gorm:"column:external_id;not null;uniqueIndex"`
	InstalledAt          time.Time `gorm:"column:installed_at;not null"`
	CreatedAt            time.Time `gorm:"column:created_at;not null"`
	UpdatedAt            time.Time `gorm:"column:updated_at;not null"`
}

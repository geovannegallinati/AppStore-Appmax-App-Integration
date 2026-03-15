package models

import "time"

type Order struct {
	ID               int64        `gorm:"column:id;primaryKey;autoIncrement"`
	InstallationID   int64        `gorm:"column:installation_id;not null"`
	Installation     Installation `gorm:"foreignKey:InstallationID"`
	AppmaxCustomerID int          `gorm:"column:appmax_customer_id;not null"`
	AppmaxOrderID    int          `gorm:"column:appmax_order_id;not null;uniqueIndex"`
	Status           string       `gorm:"column:status;not null;default:pendente"`
	PaymentMethod    string       `gorm:"column:payment_method"`
	TotalCents       int          `gorm:"column:total_cents;not null;default:0"`
	PixQRCode        string       `gorm:"column:pix_qr_code"`
	PixEMV           string       `gorm:"column:pix_emv"`
	BoletoPDFURL     string       `gorm:"column:boleto_pdf_url"`
	BoletoDigitavel  string       `gorm:"column:boleto_digitavel"`
	UpsellHash       string       `gorm:"column:upsell_hash"`
	CreatedAt        time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt        time.Time    `gorm:"column:updated_at;not null"`
}

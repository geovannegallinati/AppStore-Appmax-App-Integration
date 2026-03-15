package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateOrdersTable struct{}

func (r *CreateOrdersTable) Signature() string {
	return "20260313000002_create_orders_table"
}

func (r *CreateOrdersTable) Up() error {
	return facades.Schema().Create("orders", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.UnsignedBigInteger("installation_id")
		table.Foreign("installation_id").References("id").On("installations")
		table.Integer("appmax_customer_id")
		table.Integer("appmax_order_id")
		table.Unique("appmax_order_id")
		table.String("status", 64).Default("pendente")
		table.String("payment_method", 32).Nullable()
		table.Integer("total_cents").Default(0)
		table.Text("pix_qr_code").Nullable()
		table.Text("pix_emv").Nullable()
		table.Text("boleto_pdf_url").Nullable()
		table.Text("boleto_digitavel").Nullable()
		table.String("upsell_hash", 255).Nullable()
		table.Timestamps()
	})
}

func (r *CreateOrdersTable) Down() error {
	return facades.Schema().DropIfExists("orders")
}

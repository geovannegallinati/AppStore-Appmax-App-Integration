package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateWebhookEventsTable struct{}

func (r *CreateWebhookEventsTable) Signature() string {
	return "20260313000003_create_webhook_events_table"
}

func (r *CreateWebhookEventsTable) Up() error {
	return facades.Schema().Create("webhook_events", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("event", 64)
		table.String("event_type", 32)
		table.Integer("appmax_order_id").Nullable()
		table.Json("payload")
		table.Boolean("processed").Default(false)
		table.Timestamp("processed_at").Nullable()
		table.Text("error_message").Nullable()
		table.Timestamp("created_at").UseCurrent()
	})
}

func (r *CreateWebhookEventsTable) Down() error {
	return facades.Schema().DropIfExists("webhook_events")
}

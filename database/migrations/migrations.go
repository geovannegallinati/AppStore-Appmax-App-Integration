package migrations

import "github.com/goravel/framework/contracts/database/schema"

func All() []schema.Migration {
	return []schema.Migration{
		&CreateInstallationsTable{},
		&CreateOrdersTable{},
		&CreateWebhookEventsTable{},
	}
}

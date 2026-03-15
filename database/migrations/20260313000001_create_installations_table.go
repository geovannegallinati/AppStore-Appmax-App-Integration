package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type CreateInstallationsTable struct{}

func (r *CreateInstallationsTable) Signature() string {
	return "20260313000001_create_installations_table"
}

func (r *CreateInstallationsTable) Up() error {
	return facades.Schema().Create("installations", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("external_key", 255)
		table.String("app_id", 255)
		table.String("merchant_client_id", 512)
		table.String("merchant_client_secret", 512)
		table.Uuid("external_id")
		table.Unique("external_key")
		table.Unique("external_id")
		table.Timestamp("installed_at").UseCurrent()
		table.Timestamps()
	})
}

func (r *CreateInstallationsTable) Down() error {
	return facades.Schema().DropIfExists("installations")
}

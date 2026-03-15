package integration_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		os.Exit(0)
	}

	var err error
	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("open db: " + err.Error())
	}
	if err = testDB.Ping(); err != nil {
		panic("ping db: " + err.Error())
	}

	setupSchema()

	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func setupSchema() {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS installations (
			id                     BIGSERIAL    PRIMARY KEY,
			external_key           VARCHAR(255) NOT NULL,
			app_id                 VARCHAR(255) NOT NULL,
			merchant_client_id     VARCHAR(512) NOT NULL,
			merchant_client_secret VARCHAR(512) NOT NULL,
			external_id            UUID         NOT NULL DEFAULT gen_random_uuid(),
			installed_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			created_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at             TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			CONSTRAINT uq_int_external_key UNIQUE (external_key),
			CONSTRAINT uq_int_external_id  UNIQUE (external_id)
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id                  BIGSERIAL    PRIMARY KEY,
			installation_id     BIGINT       NOT NULL REFERENCES installations(id),
			appmax_customer_id  INTEGER      NOT NULL,
			appmax_order_id     INTEGER      NOT NULL,
			status              VARCHAR(64)  NOT NULL DEFAULT 'pendente',
			payment_method      VARCHAR(32),
			total_cents         INTEGER      NOT NULL DEFAULT 0,
			pix_qr_code         TEXT,
			pix_emv             TEXT,
			boleto_pdf_url      TEXT,
			boleto_digitavel    TEXT,
			upsell_hash         VARCHAR(255),
			created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
			CONSTRAINT uq_ord_appmax_order_id UNIQUE (appmax_order_id)
		)`,
		`CREATE TABLE IF NOT EXISTS webhook_events (
			id              BIGSERIAL    PRIMARY KEY,
			event           VARCHAR(64)  NOT NULL,
			event_type      VARCHAR(32)  NOT NULL,
			appmax_order_id INTEGER,
			payload         JSONB        NOT NULL,
			processed       BOOLEAN      NOT NULL DEFAULT FALSE,
			processed_at    TIMESTAMPTZ,
			error_message   TEXT,
			created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
		)`,
	}
	for _, stmt := range stmts {
		if _, err := testDB.Exec(stmt); err != nil {
			panic("setup schema: " + err.Error())
		}
	}
}

func truncateTables(t *testing.T) {
	t.Helper()
	_, err := testDB.Exec(`TRUNCATE webhook_events, orders, installations RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("truncate tables: %v", err)
	}
}

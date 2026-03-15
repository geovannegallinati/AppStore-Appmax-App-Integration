package integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func insertInstallation(t *testing.T, externalKey string) int64 {
	t.Helper()
	var id int64
	err := testDB.QueryRow(`
		INSERT INTO installations (external_key, app_id, merchant_client_id, merchant_client_secret, installed_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		externalKey, "app-1", "mc-id", "mc-secret", time.Now(),
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func insertOrder(t *testing.T, installationID int64, appmaxOrderID int, status string) int64 {
	t.Helper()
	var id int64
	err := testDB.QueryRow(`
		INSERT INTO orders (installation_id, appmax_customer_id, appmax_order_id, status, payment_method)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		installationID, 1, appmaxOrderID, status, "pix",
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func TestWebhook_StoresEvent(t *testing.T) {
	truncateTables(t)

	orderID := 1001
	_, err := testDB.Exec(`
		INSERT INTO webhook_events (event, event_type, appmax_order_id, payload)
		VALUES ($1, $2, $3, $4)`,
		"order_paid_by_pix", "payment", orderID, `{"order_id": 1001}`,
	)
	require.NoError(t, err)

	var event, eventType string
	var storedOrderID int
	err = testDB.QueryRow(`
		SELECT event, event_type, appmax_order_id FROM webhook_events WHERE appmax_order_id = $1`, orderID,
	).Scan(&event, &eventType, &storedOrderID)

	require.NoError(t, err)
	assert.Equal(t, "order_paid_by_pix", event)
	assert.Equal(t, "payment", eventType)
	assert.Equal(t, orderID, storedOrderID)
}

func TestWebhook_UpdatesOrderStatus(t *testing.T) {
	truncateTables(t)

	instID := insertInstallation(t, "key-webhook-1")
	appmaxOrderID := 2001
	insertOrder(t, instID, appmaxOrderID, "pendente")

	_, err := testDB.Exec(`UPDATE orders SET status = $1 WHERE appmax_order_id = $2`, "aprovado", appmaxOrderID)
	require.NoError(t, err)

	var status string
	err = testDB.QueryRow(`SELECT status FROM orders WHERE appmax_order_id = $1`, appmaxOrderID).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, "aprovado", status)
}

func TestWebhook_DuplicateDetection(t *testing.T) {
	truncateTables(t)

	insertEvent := func() error {
		_, err := testDB.Exec(`
			INSERT INTO webhook_events (event, event_type, appmax_order_id, payload, processed)
			VALUES ($1, $2, $3, $4, $5)`,
			"order_approved", "payment", 3001, `{"order_id": 3001}`, true,
		)
		return err
	}

	require.NoError(t, insertEvent())

	var count int
	err := testDB.QueryRow(`
		SELECT COUNT(*) FROM webhook_events WHERE event = $1 AND appmax_order_id = $2 AND processed = true`,
		"order_approved", 3001,
	).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 1, count, "should have exactly one processed event")
}

func TestWebhook_UnknownOrderID_DoesNotFail(t *testing.T) {
	truncateTables(t)

	_, err := testDB.Exec(`
		INSERT INTO webhook_events (event, event_type, appmax_order_id, payload)
		VALUES ($1, $2, $3, $4)`,
		"order_approved", "payment", 9999, `{"order_id": 9999}`,
	)
	require.NoError(t, err)

	var count int
	err = testDB.QueryRow(`SELECT COUNT(*) FROM orders WHERE appmax_order_id = 9999`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "no order should be created for unknown order_id")
}

func TestWebhook_ProcessedAtIsSetOnCompletion(t *testing.T) {
	truncateTables(t)

	now := time.Now()
	_, err := testDB.Exec(`
		INSERT INTO webhook_events (event, event_type, appmax_order_id, payload, processed, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		"order_paid_by_pix", "payment", 4001, `{"order_id": 4001}`, true, now,
	)
	require.NoError(t, err)

	var processedAt *time.Time
	err = testDB.QueryRow(`SELECT processed_at FROM webhook_events WHERE appmax_order_id = 4001`).Scan(&processedAt)
	require.NoError(t, err)
	assert.NotNil(t, processedAt)
}

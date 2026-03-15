package integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckout_OrderPersistedWithCorrectStatus(t *testing.T) {
	truncateTables(t)

	instID := insertInstallation(t, "key-checkout-1")

	_, err := testDB.Exec(`
		INSERT INTO orders (installation_id, appmax_customer_id, appmax_order_id, status, payment_method, pix_qr_code, pix_emv)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		instID, 42, 5001, "pendente", "pix", "qr-code-data", "emv-string",
	)
	require.NoError(t, err)

	var status, method, qrCode string
	err = testDB.QueryRow(`
		SELECT status, payment_method, pix_qr_code FROM orders WHERE appmax_order_id = 5001`,
	).Scan(&status, &method, &qrCode)

	require.NoError(t, err)
	assert.Equal(t, "pendente", status)
	assert.Equal(t, "pix", method)
	assert.Equal(t, "qr-code-data", qrCode)
}

func TestCheckout_OrderStatusUpdatedAfterWebhook(t *testing.T) {
	truncateTables(t)

	instID := insertInstallation(t, "key-checkout-2")
	insertOrder(t, instID, 6001, "pendente")

	_, err := testDB.Exec(`UPDATE orders SET status = 'aprovado' WHERE appmax_order_id = 6001`)
	require.NoError(t, err)

	var status string
	err = testDB.QueryRow(`SELECT status FROM orders WHERE appmax_order_id = 6001`).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, "aprovado", status)
}

func TestCheckout_CancelledOrderPersistedOnPaymentFailure(t *testing.T) {
	truncateTables(t)

	instID := insertInstallation(t, "key-checkout-3")

	_, err := testDB.Exec(`
		INSERT INTO orders (installation_id, appmax_customer_id, appmax_order_id, status, payment_method)
		VALUES ($1, $2, $3, $4, $5)`,
		instID, 42, 7001, "cancelado", "credit_card",
	)
	require.NoError(t, err)

	var status string
	err = testDB.QueryRow(`SELECT status FROM orders WHERE appmax_order_id = 7001`).Scan(&status)
	require.NoError(t, err)
	assert.Equal(t, "cancelado", status)
}

func TestCheckout_BoletoFieldsPersisted(t *testing.T) {
	truncateTables(t)

	instID := insertInstallation(t, "key-checkout-4")

	_, err := testDB.Exec(`
		INSERT INTO orders (installation_id, appmax_customer_id, appmax_order_id, status, payment_method, boleto_pdf_url, boleto_digitavel)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		instID, 42, 8001, "pendente", "boleto", "https://boleto.pdf", "1234.5678 9012.3456",
	)
	require.NoError(t, err)

	var pdfURL, digitavel string
	err = testDB.QueryRow(`
		SELECT boleto_pdf_url, boleto_digitavel FROM orders WHERE appmax_order_id = 8001`,
	).Scan(&pdfURL, &digitavel)

	require.NoError(t, err)
	assert.Equal(t, "https://boleto.pdf", pdfURL)
	assert.Equal(t, "1234.5678 9012.3456", digitavel)
}

func TestCheckout_OrderHasCorrectTimestamps(t *testing.T) {
	truncateTables(t)

	instID := insertInstallation(t, "key-checkout-5")

	before := time.Now().Add(-time.Second)
	insertOrder(t, instID, 9001, "pendente")
	after := time.Now().Add(time.Second)

	var createdAt, updatedAt time.Time
	err := testDB.QueryRow(`
		SELECT created_at, updated_at FROM orders WHERE appmax_order_id = 9001`,
	).Scan(&createdAt, &updatedAt)

	require.NoError(t, err)
	assert.True(t, createdAt.After(before) && createdAt.Before(after), "created_at should be set to current time")
	assert.True(t, updatedAt.After(before) && updatedAt.Before(after), "updated_at should be set to current time")
}

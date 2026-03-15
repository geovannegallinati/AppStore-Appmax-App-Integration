package integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallation_CreateAndRetrieve(t *testing.T) {
	truncateTables(t)

	_, err := testDB.Exec(`
		INSERT INTO installations (external_key, app_id, merchant_client_id, merchant_client_secret, installed_at)
		VALUES ($1, $2, $3, $4, $5)`,
		"key-abc", "app-1", "mc-id", "mc-secret", time.Now(),
	)
	require.NoError(t, err)

	var externalID, externalKey string
	err = testDB.QueryRow(
		`SELECT external_id, external_key FROM installations WHERE external_key = $1`, "key-abc",
	).Scan(&externalID, &externalKey)

	require.NoError(t, err)
	assert.Equal(t, "key-abc", externalKey)
	assert.NotEmpty(t, externalID, "external_id should be auto-generated UUID")
}

func TestInstallation_UniqueExternalKey(t *testing.T) {
	truncateTables(t)

	insert := func() error {
		_, err := testDB.Exec(`
			INSERT INTO installations (external_key, app_id, merchant_client_id, merchant_client_secret, installed_at)
			VALUES ($1, $2, $3, $4, $5)`,
			"key-dup", "app-1", "mc-id", "mc-secret", time.Now(),
		)
		return err
	}

	require.NoError(t, insert())
	err := insert()
	require.Error(t, err, "second insert with same external_key should fail")
}

func TestInstallation_UpdateMerchantCreds(t *testing.T) {
	truncateTables(t)

	_, err := testDB.Exec(`
		INSERT INTO installations (external_key, app_id, merchant_client_id, merchant_client_secret, installed_at)
		VALUES ($1, $2, $3, $4, $5)`,
		"key-update", "app-1", "old-id", "old-secret", time.Now(),
	)
	require.NoError(t, err)

	_, err = testDB.Exec(`
		UPDATE installations
		SET merchant_client_id = $1, merchant_client_secret = $2, updated_at = NOW()
		WHERE external_key = $3`,
		"new-id", "new-secret", "key-update",
	)
	require.NoError(t, err)

	var clientID string
	err = testDB.QueryRow(
		`SELECT merchant_client_id FROM installations WHERE external_key = $1`, "key-update",
	).Scan(&clientID)

	require.NoError(t, err)
	assert.Equal(t, "new-id", clientID)
}

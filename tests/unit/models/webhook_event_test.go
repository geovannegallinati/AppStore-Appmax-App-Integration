package models_test

import (
	"testing"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONMap_ValueAndScan(t *testing.T) {
	payload := models.JSONMap{"event": "order_paid", "order_id": 123}

	value, err := payload.Value()
	require.NoError(t, err)

	var scanned models.JSONMap
	require.NoError(t, scanned.Scan(value))
	assert.Equal(t, "order_paid", scanned["event"])
	assert.EqualValues(t, 123, scanned["order_id"])
}

func TestJSONMap_ScanUnsupportedType(t *testing.T) {
	var scanned models.JSONMap

	err := scanned.Scan(123)

	require.Error(t, err)
	assert.ErrorContains(t, err, "unsupported type")
}

package requests_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/http/requests"
)

func TestWebhookDataRequest_OrderIDFromNumber(t *testing.T) {
	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(`{"order_id":123}`), &data)

	require.NoError(t, err)
	require.NotNil(t, data.OrderID.Ptr())
	assert.Equal(t, 123, *data.OrderID.Ptr())
}

func TestWebhookDataRequest_OrderIDFromString(t *testing.T) {
	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(`{"order_id":"456"}`), &data)

	require.NoError(t, err)
	require.NotNil(t, data.OrderID.Ptr())
	assert.Equal(t, 456, *data.OrderID.Ptr())
}

func TestWebhookDataRequest_OrderIDNull(t *testing.T) {
	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(`{"order_id":null}`), &data)

	require.NoError(t, err)
	assert.Nil(t, data.OrderID.Ptr())
}

func TestWebhookDataRequest_OrderIDInvalid(t *testing.T) {
	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(`{"order_id":"abc"}`), &data)

	require.Error(t, err)
}

package requests_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/requests"
)

func TestWebhookDataRequest_OptionalInt_FromNumber(t *testing.T) {
	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(`{"order_id":123}`), &data)

	require.NoError(t, err)
	require.NotNil(t, data.OrderID.Ptr())
	assert.Equal(t, 123, *data.OrderID.Ptr())
}

func TestWebhookDataRequest_OptionalInt_FromString(t *testing.T) {
	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(`{"order_id":"456"}`), &data)

	require.NoError(t, err)
	require.NotNil(t, data.OrderID.Ptr())
	assert.Equal(t, 456, *data.OrderID.Ptr())
}

func TestWebhookDataRequest_OptionalInt_Null(t *testing.T) {
	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(`{"order_id":null}`), &data)

	require.NoError(t, err)
	assert.Nil(t, data.OrderID.Ptr())
}

func TestWebhookDataRequest_OptionalInt_Invalid(t *testing.T) {
	var data requests.WebhookDataRequest
	err := json.Unmarshal([]byte(`{"order_id":"abc"}`), &data)

	require.Error(t, err)
}

func TestExtractOrderID_Standard_UsesDataId(t *testing.T) {
	payload := `{"id":42,"customer_id":7,"total_products":265,"status":"aprovado","payment_type":"CreditCard","total":267.48}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	orderID := data.ExtractOrderID()
	require.NotNil(t, orderID)
	assert.Equal(t, 42, *orderID)
}

func TestExtractOrderID_StandardWithMeta_UsesDataId(t *testing.T) {
	payload := `{"id":42,"customer_id":7,"total_products":265,"status":"aprovado","payment_type":"CreditCard","total":267.48,"meta":[]}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	orderID := data.ExtractOrderID()
	require.NotNil(t, orderID)
	assert.Equal(t, 42, *orderID)
}

func TestExtractOrderID_TwoLevelFlat_UsesOrderId(t *testing.T) {
	payload := `{"order_id":99,"order_total_products":369,"order_status":"aprovado","customer_id":7,"customer_firstname":"Leandro"}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	orderID := data.ExtractOrderID()
	require.NotNil(t, orderID)
	assert.Equal(t, 99, *orderID)
}

func TestExtractOrderID_CustomContent_UsesOrderId(t *testing.T) {
	payload := `{"order_id":12844,"foo":"bar"}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	orderID := data.ExtractOrderID()
	require.NotNil(t, orderID)
	assert.Equal(t, 12844, *orderID)
}

func TestExtractOrderID_OldLegacy_UsesOrderId(t *testing.T) {
	payload := `{"order_id":12764}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	orderID := data.ExtractOrderID()
	require.NotNil(t, orderID)
	assert.Equal(t, 12764, *orderID)
}

func TestExtractOrderID_CustomerEvent_ReturnsNil(t *testing.T) {
	payload := `{"id":1,"site_id":1470,"firstname":"Laura","lastname":"Montenegro","email":"maximo33@gmail.com"}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	orderID := data.ExtractOrderID()
	assert.Nil(t, orderID, "customer events have no customer_id in data, so data.id must not be used as order ID")
}

func TestExtractOrderID_SubscriptionClientEvent_ReturnsNil(t *testing.T) {
	payload := `{"id":1,"site_id":1470,"firstname":"Noeli","lastname":"Guerra","subscription":{"id":null}}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	orderID := data.ExtractOrderID()
	assert.Nil(t, orderID, "SubscriptionCancellationEvent uses customer structure — data.id is customer ID not order ID")
}

func TestDetectModel_Standard(t *testing.T) {
	payload := `{"id":42,"customer_id":7,"total_products":265,"status":"aprovado","payment_type":"CreditCard"}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	assert.Equal(t, "standard", data.DetectModel(""))
}

func TestDetectModel_StandardWithMeta(t *testing.T) {
	payload := `{"id":42,"customer_id":7,"total_products":265,"status":"aprovado","payment_type":"CreditCard","meta":[]}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	assert.Equal(t, "standard_meta", data.DetectModel(""))
}

func TestDetectModel_TwoLevelFlat(t *testing.T) {
	payload := `{"order_id":99,"order_total_products":369,"order_status":"aprovado","customer_id":7,"customer_firstname":"Leandro"}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	assert.Equal(t, "two_level_flat", data.DetectModel(""))
}

func TestDetectModel_CustomContent(t *testing.T) {
	payload := `{"order_id":12844,"foo":"bar"}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	assert.Equal(t, "custom_content", data.DetectModel(""))
}

func TestDetectModel_OldLegacy(t *testing.T) {
	payload := `{"order_id":12764}`

	var data requests.WebhookDataRequest
	require.NoError(t, json.Unmarshal([]byte(payload), &data))

	assert.Equal(t, "old_legacy", data.DetectModel("order"))
}

//go:build appmax_live

package appmax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

func TestOrderCreation(t *testing.T) {
	cid := ensureCustomerID(t)
	id := createOrderForCustomer(t, cid, "PROD-001", "Produto E2E", 5000, "digital")
	assert.Greater(t, id, 0)
	t.Logf("Order created: id=%d", id)
}

func TestGetOrder(t *testing.T) {
	cid := ensureCustomerID(t)
	oid := createOrderForCustomer(t, cid, "PROD-GET", "Produto GetOrder", 5300, "digital")
	result, err := appmaxSvc.GetOrder(testCtx(t), testInst, oid)
	require.NoError(t, err)
	assert.Equal(t, oid, result.OrderID)
	assert.NotEmpty(t, result.Status)
	t.Logf("Order %d: status=%s total_paid=%d customer=%s",
		result.OrderID, result.Status, result.TotalPaid, result.CustomerName)
}

func TestAddOrderTracking(t *testing.T) {
	cid := ensureCustomerID(t)
	oid := createOrderForCustomer(t, cid, "PROD-TRACK", "Produto Tracking", 5400, "physical")

	err := appmaxSvc.AddTracking(testCtx(t), testInst, services.TrackingInput{
		OrderID:              oid,
		ShippingTrackingCode: "BR123456789XX",
	})
	require.NoError(t, err)
	t.Logf("Tracking code added to order %d", oid)
}

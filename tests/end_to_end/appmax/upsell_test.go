//go:build appmax_live

package appmax

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

func TestUpsell(t *testing.T) {
	hash := ensureUpsellHash(t)

	result, err := appmaxSvc.Upsell(testCtx(t), testInst, services.UpsellInput{
		UpsellHash:    hash,
		ProductsValue: 2000,
		Products: []services.Product{
			{SKU: "PROD-UP", Name: "Upsell Product", Quantity: 1, UnitValue: 2000, Type: "digital"},
		},
	})
	require.NoError(t, err)
	t.Logf("Upsell order: id=%d status=%s", result.OrderID, result.Status)
}

//go:build appmax_live

package appmax

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

func TestTotalRefund(t *testing.T) {
	oid, _, _ := createApprovedCreditCardPayment(t)

	err := appmaxSvc.Refund(testCtx(t), testInst, services.RefundInput{
		OrderID: oid,
		Type:    "total",
	})
	require.NoError(t, err)
	t.Logf("Total refund completed for order %d", oid)
}

func TestPartialRefund(t *testing.T) {
	partialOrderID, _, _ := createApprovedCreditCardPayment(t)

	refundErr := appmaxSvc.Refund(testCtx(t), testInst, services.RefundInput{
		OrderID: partialOrderID,
		Type:    "partial",
		Value:   5000,
	})
	require.NoError(t, refundErr)
	t.Logf("Partial refund (R$50.00) completed for order %d", partialOrderID)
}

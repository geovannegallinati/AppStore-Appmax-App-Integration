//go:build appmax_live

package appmax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

func TestCreditCardPaymentApproved(t *testing.T) {
	cid := ensureCustomerID(t)
	oid := createOrderForCustomer(t, cid, "PROD-APPROVED", "Produto E2E Aprovado", 5000, "digital")

	result, err := appmaxSvc.CreditCard(testCtx(t), testInst, services.CreditCardInput{
		OrderID:              oid,
		CustomerID:           cid,
		Number:               cardSuccess,
		CVV:                  "123",
		ExpirationMonth:      "12",
		ExpirationYear:       "28",
		HolderName:           "TESTE E2E",
		HolderDocumentNumber: docCPF,
		Installments:         1,
		SoftDescriptor:       "E2ETEST",
	})
	require.NoError(t, err)
	t.Logf("Credit card payment: id=%d status=%s upsell_hash=%s",
		result.PaymentID, result.Status, result.UpsellHash)

	if result.UpsellHash != "" {
		fixtureMu.Lock()
		upsellHash = result.UpsellHash
		fixtureMu.Unlock()
	}
}

func TestCreditCardPaymentDeclined(t *testing.T) {
	cid := ensureCustomerID(t)

	failOrderID, orderErr := appmaxSvc.CreateOrder(testCtx(t), testInst, services.OrderInput{
		CustomerID:    cid,
		DiscountValue: 0,
		ShippingValue: 0,
		Products: []services.Product{
			{SKU: "PROD-FAIL", Name: "Produto Falha", Quantity: 1, UnitValue: 1000, Type: "digital"},
		},
	})
	require.NoError(t, orderErr)

	_, err := appmaxSvc.CreditCard(testCtx(t), testInst, services.CreditCardInput{
		OrderID:              failOrderID,
		CustomerID:           cid,
		Number:               cardFail,
		CVV:                  "123",
		ExpirationMonth:      "12",
		ExpirationYear:       "28",
		HolderName:           "TESTE E2E",
		HolderDocumentNumber: docCPF,
		Installments:         1,
	})
	assert.Error(t, err, "failure card must return an error")
	t.Logf("Card decline (expected): %v", err)
}

func TestPixPayment(t *testing.T) {
	cid := ensureCustomerID(t)

	pixOrderID, orderErr := appmaxSvc.CreateOrder(testCtx(t), testInst, services.OrderInput{
		CustomerID:    cid,
		DiscountValue: 0,
		ShippingValue: 0,
		Products: []services.Product{
			{SKU: "PROD-PIX", Name: "Produto Pix", Quantity: 1, UnitValue: 3000, Type: "digital"},
		},
	})
	require.NoError(t, orderErr)

	result, err := appmaxSvc.Pix(testCtx(t), testInst, services.PixInput{
		OrderID:        pixOrderID,
		DocumentNumber: docCPF,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.QRCode, "QR code must not be empty")
	t.Logf("Pix order=%d qr_code=%s...", pixOrderID, truncate(result.QRCode, 40))
}

func TestBoletoPayment(t *testing.T) {
	cid := ensureCustomerID(t)

	boletoOrderID, orderErr := appmaxSvc.CreateOrder(testCtx(t), testInst, services.OrderInput{
		CustomerID:    cid,
		DiscountValue: 0,
		ShippingValue: 0,
		Products: []services.Product{
			{SKU: "PROD-BOL", Name: "Produto Boleto", Quantity: 1, UnitValue: 7500, Type: "physical"},
		},
	})
	require.NoError(t, orderErr)

	result, err := appmaxSvc.Boleto(testCtx(t), testInst, services.BoletoInput{
		OrderID:        boletoOrderID,
		DocumentNumber: docCPF,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result.PDFURL, "boleto PDF URL must not be empty")
	t.Logf("Boleto order=%d pdf_url=%s", boletoOrderID, result.PDFURL)
}

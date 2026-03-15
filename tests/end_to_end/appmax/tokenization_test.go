//go:build appmax_live

package appmax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

func TestInstallmentCalculation(t *testing.T) {
	items, err := appmaxSvc.Installments(testCtx(t), testInst, services.InstallmentsInput{
		Installments: 12,
		TotalValue:   10000,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, items, "installments list must not be empty")
	t.Logf("Installment options returned: %d", len(items))
	for _, item := range items {
		t.Logf("  %dx = R$%.2f (total R$%.2f)", item.Installments, item.Value/100, item.TotalValue/100)
	}
}

func TestCardTokenization(t *testing.T) {
	token, err := appmaxSvc.Tokenize(testCtx(t), testInst, services.TokenizeInput{
		Number:          cardSuccess,
		CVV:             "123",
		ExpirationMonth: "12",
		ExpirationYear:  "28",
		HolderName:      "TESTE E2E",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, token, "token must not be empty")
	t.Logf("Card token: %s", token)
}

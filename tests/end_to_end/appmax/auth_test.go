//go:build appmax_live

package appmax

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppTokenIsValid(t *testing.T) {
	require.NotEmpty(t, creds.AppToken, "app token must be set")
	t.Logf("App token: %s...", creds.AppToken[:20])
}

func TestMerchantTokenIsValid(t *testing.T) {
	require.NotEmpty(t, creds.MerchantToken, "merchant token must be set")
	t.Logf("Merchant token: %s...", creds.MerchantToken[:20])
}

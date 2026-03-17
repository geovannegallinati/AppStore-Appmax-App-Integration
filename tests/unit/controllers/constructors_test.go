package controllers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/controllers"
)

func TestControllerConstructors_RejectNilDependency(t *testing.T) {
	merchantAuthController, err := controllers.NewMerchantAuthController(nil)
	require.Error(t, err)
	assert.Nil(t, merchantAuthController)

	checkoutController, err := controllers.NewCheckoutController(nil)
	require.Error(t, err)
	assert.Nil(t, checkoutController)

	webhookController, err := controllers.NewWebhookController(nil)
	require.Error(t, err)
	assert.Nil(t, webhookController)

	installController, err := controllers.NewInstallController(nil, nil, "", "", "", "")
	require.Error(t, err)
	assert.Nil(t, installController)
}

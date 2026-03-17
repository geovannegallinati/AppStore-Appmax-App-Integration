package controllers_test

import (
	"context"
	"testing"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/controllers"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/models"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopAppmaxService struct{}

func (noopAppmaxService) Authorize(context.Context, string, string, string) (string, error) {
	return "", nil
}
func (noopAppmaxService) GenerateMerchantCreds(context.Context, string) (string, string, error) {
	return "", "", nil
}
func (noopAppmaxService) CreateOrUpdateCustomer(context.Context, *models.Installation, services.CustomerInput) (int, error) {
	return 0, nil
}
func (noopAppmaxService) CreateOrder(context.Context, *models.Installation, services.OrderInput) (int, error) {
	return 0, nil
}
func (noopAppmaxService) CreditCard(context.Context, *models.Installation, services.CreditCardInput) (services.CreditCardResult, error) {
	return services.CreditCardResult{}, nil
}
func (noopAppmaxService) Pix(context.Context, *models.Installation, services.PixInput) (services.PixResult, error) {
	return services.PixResult{}, nil
}
func (noopAppmaxService) Boleto(context.Context, *models.Installation, services.BoletoInput) (services.BoletoResult, error) {
	return services.BoletoResult{}, nil
}
func (noopAppmaxService) Installments(context.Context, *models.Installation, services.InstallmentsInput) ([]services.AppmaxInstallmentItem, error) {
	return nil, nil
}
func (noopAppmaxService) Refund(context.Context, *models.Installation, services.RefundInput) error {
	return nil
}
func (noopAppmaxService) GetOrder(context.Context, *models.Installation, int) (services.GetOrderResult, error) {
	return services.GetOrderResult{}, nil
}
func (noopAppmaxService) AddTracking(context.Context, *models.Installation, services.TrackingInput) error {
	return nil
}
func (noopAppmaxService) Upsell(context.Context, *models.Installation, services.UpsellInput) (services.UpsellResult, error) {
	return services.UpsellResult{}, nil
}
func (noopAppmaxService) Tokenize(context.Context, *models.Installation, services.TokenizeInput) (string, error) {
	return "", nil
}
func (noopAppmaxService) VerifyMerchantCreds(context.Context, string, string) error { return nil }

type noopInstallService struct{}

func (noopInstallService) Upsert(context.Context, services.UpsertInstallationInput) (*models.Installation, bool, error) {
	return &models.Installation{}, true, nil
}

type noopTokenManager struct{}

func (noopTokenManager) AppToken(context.Context) (string, error) {
	return "app-token", nil
}

func (noopTokenManager) MerchantToken(context.Context, *models.Installation) (string, error) {
	return "merchant-token", nil
}

type noopCheckoutService struct{}

func (noopCheckoutService) ProcessCreditCard(context.Context, *models.Installation, services.CheckoutCreditCardInput) (services.CheckoutCreditCardResult, error) {
	return services.CheckoutCreditCardResult{}, nil
}
func (noopCheckoutService) ProcessPix(context.Context, *models.Installation, services.CheckoutPixInput) (services.CheckoutPixResult, error) {
	return services.CheckoutPixResult{}, nil
}
func (noopCheckoutService) ProcessBoleto(context.Context, *models.Installation, services.CheckoutBoletoInput) (services.CheckoutBoletoResult, error) {
	return services.CheckoutBoletoResult{}, nil
}
func (noopCheckoutService) GetOrderStatus(context.Context, *models.Installation, int) (string, error) {
	return "", nil
}
func (noopCheckoutService) GetInstallments(context.Context, *models.Installation, services.InstallmentsInput) ([]services.AppmaxInstallmentItem, error) {
	return nil, nil
}
func (noopCheckoutService) ProcessRefund(context.Context, *models.Installation, services.RefundInput) error {
	return nil
}
func (noopCheckoutService) Tokenize(context.Context, *models.Installation, services.TokenizeInput) (string, error) {
	return "", nil
}
func (noopCheckoutService) AddTracking(context.Context, *models.Installation, services.TrackingInput) error {
	return nil
}
func (noopCheckoutService) CreateCustomerAndOrder(context.Context, *models.Installation, services.CustomerInput, services.OrderInput) (services.CheckoutCreateOrderResult, error) {
	return services.CheckoutCreateOrderResult{}, nil
}
func (noopCheckoutService) ProcessUpsell(context.Context, *models.Installation, services.UpsellInput) (services.UpsellResult, error) {
	return services.UpsellResult{}, nil
}

type noopWebhookService struct{}

func (noopWebhookService) Handle(context.Context, services.WebhookInput) (services.WebhookResult, error) {
	return services.WebhookResult{}, nil
}

func TestControllerConstructors_Success(t *testing.T) {
	merchantAuthController, err := controllers.NewMerchantAuthController(noopTokenManager{})
	require.NoError(t, err)
	assert.NotNil(t, merchantAuthController)

	checkoutController, err := controllers.NewCheckoutController(noopCheckoutService{})
	require.NoError(t, err)
	assert.NotNil(t, checkoutController)

	webhookController, err := controllers.NewWebhookController(noopWebhookService{}, "https://admin.appmax.com.br", "https://app.example.com")
	require.NoError(t, err)
	assert.NotNil(t, webhookController)

	installController, err := controllers.NewInstallController(noopAppmaxService{}, noopInstallService{}, "https://admin.appmax.com.br", "https://app.example.com", "test-app-uuid", "123")
	require.NoError(t, err)
	assert.NotNil(t, installController)
}

func TestInstallControllerConstructor_RejectsEmptyAdminURL(t *testing.T) {
	installController, err := controllers.NewInstallController(noopAppmaxService{}, noopInstallService{}, "", "https://app.example.com", "test-app-uuid", "123")

	require.Error(t, err)
	assert.Nil(t, installController)
	assert.ErrorIs(t, err, controllers.ErrInvalidConfig)
}

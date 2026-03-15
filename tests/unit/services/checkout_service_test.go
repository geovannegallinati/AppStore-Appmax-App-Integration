package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
	"github.com/geovanne-gallinati/AppStoreAppDemo/tests/unit/mocks"
)

type mockAppmaxService struct {
	authorizeFunc              func(ctx context.Context, appID, externalKey, callbackURL string) (string, error)
	generateMerchantCredsFunc  func(ctx context.Context, hash string) (string, string, error)
	createOrUpdateCustomerFunc func(ctx context.Context, inst *models.Installation, input services.CustomerInput) (int, error)
	createOrderFunc            func(ctx context.Context, inst *models.Installation, input services.OrderInput) (int, error)
	creditCardFunc             func(ctx context.Context, inst *models.Installation, input services.CreditCardInput) (services.CreditCardResult, error)
	pixFunc                    func(ctx context.Context, inst *models.Installation, input services.PixInput) (services.PixResult, error)
	boletoFunc                 func(ctx context.Context, inst *models.Installation, input services.BoletoInput) (services.BoletoResult, error)
	installmentsFunc           func(ctx context.Context, inst *models.Installation, input services.InstallmentsInput) ([]services.AppmaxInstallmentItem, error)
	refundFunc                 func(ctx context.Context, inst *models.Installation, input services.RefundInput) error
	getOrderFunc               func(ctx context.Context, inst *models.Installation, orderID int) (services.GetOrderResult, error)
	addTrackingFunc            func(ctx context.Context, inst *models.Installation, input services.TrackingInput) error
	upsellFunc                 func(ctx context.Context, inst *models.Installation, input services.UpsellInput) (services.UpsellResult, error)
	tokenizeFunc               func(ctx context.Context, inst *models.Installation, input services.TokenizeInput) (string, error)
}

func (m *mockAppmaxService) Authorize(ctx context.Context, appID, externalKey, callbackURL string) (string, error) {
	return m.authorizeFunc(ctx, appID, externalKey, callbackURL)
}
func (m *mockAppmaxService) GenerateMerchantCreds(ctx context.Context, hash string) (string, string, error) {
	return m.generateMerchantCredsFunc(ctx, hash)
}
func (m *mockAppmaxService) CreateOrUpdateCustomer(ctx context.Context, inst *models.Installation, input services.CustomerInput) (int, error) {
	return m.createOrUpdateCustomerFunc(ctx, inst, input)
}
func (m *mockAppmaxService) CreateOrder(ctx context.Context, inst *models.Installation, input services.OrderInput) (int, error) {
	return m.createOrderFunc(ctx, inst, input)
}
func (m *mockAppmaxService) CreditCard(ctx context.Context, inst *models.Installation, input services.CreditCardInput) (services.CreditCardResult, error) {
	return m.creditCardFunc(ctx, inst, input)
}
func (m *mockAppmaxService) Pix(ctx context.Context, inst *models.Installation, input services.PixInput) (services.PixResult, error) {
	return m.pixFunc(ctx, inst, input)
}
func (m *mockAppmaxService) Boleto(ctx context.Context, inst *models.Installation, input services.BoletoInput) (services.BoletoResult, error) {
	return m.boletoFunc(ctx, inst, input)
}
func (m *mockAppmaxService) Installments(ctx context.Context, inst *models.Installation, input services.InstallmentsInput) ([]services.AppmaxInstallmentItem, error) {
	return m.installmentsFunc(ctx, inst, input)
}
func (m *mockAppmaxService) Refund(ctx context.Context, inst *models.Installation, input services.RefundInput) error {
	return m.refundFunc(ctx, inst, input)
}
func (m *mockAppmaxService) GetOrder(ctx context.Context, inst *models.Installation, orderID int) (services.GetOrderResult, error) {
	return m.getOrderFunc(ctx, inst, orderID)
}
func (m *mockAppmaxService) AddTracking(ctx context.Context, inst *models.Installation, input services.TrackingInput) error {
	return m.addTrackingFunc(ctx, inst, input)
}
func (m *mockAppmaxService) Upsell(ctx context.Context, inst *models.Installation, input services.UpsellInput) (services.UpsellResult, error) {
	return m.upsellFunc(ctx, inst, input)
}
func (m *mockAppmaxService) Tokenize(ctx context.Context, inst *models.Installation, input services.TokenizeInput) (string, error) {
	return m.tokenizeFunc(ctx, inst, input)
}
func (m *mockAppmaxService) VerifyMerchantCreds(context.Context, string, string) error { return nil }

var checkoutTestInst = &models.Installation{
	ID:                   1,
	ExternalKey:          "test-key",
	MerchantClientID:     "mc-id",
	MerchantClientSecret: "mc-secret",
}

func testCheckoutInput() services.CheckoutCreditCardInput {
	return services.CheckoutCreditCardInput{
		Customer: services.CustomerInput{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
		Order:    services.OrderInput{ProductsValue: 10000, Products: []services.Product{{SKU: "P1", Name: "Product 1", Quantity: 1}}},
		Payment:  services.CreditCardInput{Number: "4000000000000010", Installments: 1},
	}
}

func noopOrderRepo() *mocks.MockOrderRepository {
	return &mocks.MockOrderRepository{
		CreateFunc: func(_ context.Context, _ *models.Order) error { return nil },
		SaveFunc:   func(_ context.Context, _ *models.Order) error { return nil },
		FindByAppmaxOrderIDFunc: func(_ context.Context, _ int) (*models.Order, error) {
			return nil, nil
		},
		FindByAppmaxOrderIDAndInstallationFunc: func(_ context.Context, _ int, _ int64) (*models.Order, error) {
			return nil, nil
		},
	}
}

func mustCheckoutService(t *testing.T, appmaxSvc services.AppmaxService, orderRepo contracts.OrderRepository) services.CheckoutService {
	t.Helper()

	svc, err := services.NewCheckoutService(appmaxSvc, orderRepo)
	require.NoError(t, err)

	return svc
}

func TestCheckoutService_ProcessCreditCard_Approved(t *testing.T) {
	mock := &mockAppmaxService{
		createOrUpdateCustomerFunc: func(_ context.Context, _ *models.Installation, _ services.CustomerInput) (int, error) {
			return 42, nil
		},
		createOrderFunc: func(_ context.Context, _ *models.Installation, _ services.OrderInput) (int, error) {
			return 99, nil
		},
		creditCardFunc: func(_ context.Context, _ *models.Installation, _ services.CreditCardInput) (services.CreditCardResult, error) {
			return services.CreditCardResult{PaymentID: 1, Status: "autorizado"}, nil
		},
	}

	svc := mustCheckoutService(t, mock, noopOrderRepo())
	result, err := svc.ProcessCreditCard(context.Background(), checkoutTestInst, testCheckoutInput())

	require.NoError(t, err)
	assert.Equal(t, 99, result.OrderID)
	assert.Equal(t, "autorizado", result.Status)
}

func TestCheckoutService_ProcessCreditCard_PaymentDeclined(t *testing.T) {
	mock := &mockAppmaxService{
		createOrUpdateCustomerFunc: func(_ context.Context, _ *models.Installation, _ services.CustomerInput) (int, error) {
			return 42, nil
		},
		createOrderFunc: func(_ context.Context, _ *models.Installation, _ services.OrderInput) (int, error) {
			return 99, nil
		},
		creditCardFunc: func(_ context.Context, _ *models.Installation, _ services.CreditCardInput) (services.CreditCardResult, error) {
			return services.CreditCardResult{}, services.ErrPaymentDeclined
		},
	}

	svc := mustCheckoutService(t, mock, noopOrderRepo())
	_, err := svc.ProcessCreditCard(context.Background(), checkoutTestInst, testCheckoutInput())

	require.Error(t, err)
	assert.True(t, errors.Is(err, services.ErrPaymentDeclined))
}

func TestCheckoutService_ProcessCreditCard_CustomerCreationFails(t *testing.T) {
	customerErr := errors.New("customer API down")
	mock := &mockAppmaxService{
		createOrUpdateCustomerFunc: func(_ context.Context, _ *models.Installation, _ services.CustomerInput) (int, error) {
			return 0, customerErr
		},
	}

	svc := mustCheckoutService(t, mock, noopOrderRepo())
	_, err := svc.ProcessCreditCard(context.Background(), checkoutTestInst, testCheckoutInput())

	require.Error(t, err)
	assert.ErrorContains(t, err, "customer API down")
}

func TestCheckoutService_ProcessPix_Success(t *testing.T) {
	mock := &mockAppmaxService{
		createOrUpdateCustomerFunc: func(_ context.Context, _ *models.Installation, _ services.CustomerInput) (int, error) {
			return 42, nil
		},
		createOrderFunc: func(_ context.Context, _ *models.Installation, _ services.OrderInput) (int, error) {
			return 99, nil
		},
		pixFunc: func(_ context.Context, _ *models.Installation, _ services.PixInput) (services.PixResult, error) {
			return services.PixResult{QRCode: "qr-data", EMV: "emv-string"}, nil
		},
	}

	svc := mustCheckoutService(t, mock, noopOrderRepo())
	result, err := svc.ProcessPix(context.Background(), checkoutTestInst, services.CheckoutPixInput{
		Customer:       services.CustomerInput{FirstName: "John"},
		Order:          services.OrderInput{ProductsValue: 10000},
		DocumentNumber: "12345678901",
	})

	require.NoError(t, err)
	assert.Equal(t, 99, result.OrderID)
	assert.Equal(t, "qr-data", result.QRCode)
	assert.Equal(t, "emv-string", result.EMV)
}

func TestCheckoutService_ProcessBoleto_Success(t *testing.T) {
	mock := &mockAppmaxService{
		createOrUpdateCustomerFunc: func(_ context.Context, _ *models.Installation, _ services.CustomerInput) (int, error) {
			return 42, nil
		},
		createOrderFunc: func(_ context.Context, _ *models.Installation, _ services.OrderInput) (int, error) {
			return 99, nil
		},
		boletoFunc: func(_ context.Context, _ *models.Installation, _ services.BoletoInput) (services.BoletoResult, error) {
			return services.BoletoResult{PDFURL: "https://boleto", Digitavel: "34191.79001"}, nil
		},
	}

	svc := mustCheckoutService(t, mock, noopOrderRepo())
	result, err := svc.ProcessBoleto(context.Background(), checkoutTestInst, services.CheckoutBoletoInput{
		Customer:       services.CustomerInput{FirstName: "John"},
		Order:          services.OrderInput{ProductsValue: 10000},
		DocumentNumber: "12345678901",
	})

	require.NoError(t, err)
	assert.Equal(t, 99, result.OrderID)
	assert.Equal(t, "https://boleto", result.PDFURL)
	assert.Equal(t, "34191.79001", result.Digitavel)
}

func TestCheckoutService_ProcessBoleto_Error(t *testing.T) {
	mock := &mockAppmaxService{
		createOrUpdateCustomerFunc: func(_ context.Context, _ *models.Installation, _ services.CustomerInput) (int, error) {
			return 42, nil
		},
		createOrderFunc: func(_ context.Context, _ *models.Installation, _ services.OrderInput) (int, error) {
			return 99, nil
		},
		boletoFunc: func(_ context.Context, _ *models.Installation, _ services.BoletoInput) (services.BoletoResult, error) {
			return services.BoletoResult{}, errors.New("gateway unavailable")
		},
	}

	svc := mustCheckoutService(t, mock, noopOrderRepo())
	_, err := svc.ProcessBoleto(context.Background(), checkoutTestInst, services.CheckoutBoletoInput{
		Customer:       services.CustomerInput{FirstName: "John"},
		Order:          services.OrderInput{ProductsValue: 10000},
		DocumentNumber: "12345678901",
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "gateway unavailable")
}

func TestCheckoutService_GetInstallments(t *testing.T) {
	expected := []services.AppmaxInstallmentItem{
		{Installments: 1, Value: 10000, TotalValue: 10000},
		{Installments: 2, Value: 5100, TotalValue: 10200},
	}
	mock := &mockAppmaxService{
		installmentsFunc: func(_ context.Context, _ *models.Installation, _ services.InstallmentsInput) ([]services.AppmaxInstallmentItem, error) {
			return expected, nil
		},
	}

	svc := mustCheckoutService(t, mock, noopOrderRepo())
	items, err := svc.GetInstallments(context.Background(), checkoutTestInst, services.InstallmentsInput{Installments: 12, TotalValue: 10000})

	require.NoError(t, err)
	assert.Equal(t, expected, items)
}

func TestCheckoutService_ProcessRefund_Success(t *testing.T) {
	mock := &mockAppmaxService{
		refundFunc: func(_ context.Context, _ *models.Installation, _ services.RefundInput) error {
			return nil
		},
	}

	svc := mustCheckoutService(t, mock, noopOrderRepo())
	err := svc.ProcessRefund(context.Background(), checkoutTestInst, services.RefundInput{OrderID: 99, Type: "full"})

	require.NoError(t, err)
}

func TestCheckoutService_GetOrderStatus_Found(t *testing.T) {
	orderRepo := &mocks.MockOrderRepository{
		FindByAppmaxOrderIDAndInstallationFunc: func(_ context.Context, appmaxOrderID int, installationID int64) (*models.Order, error) {
			assert.Equal(t, 99, appmaxOrderID)
			assert.Equal(t, int64(1), installationID)
			return &models.Order{AppmaxOrderID: 99, Status: "aprovado"}, nil
		},
	}

	svc := mustCheckoutService(t, &mockAppmaxService{}, orderRepo)
	status, err := svc.GetOrderStatus(context.Background(), checkoutTestInst, 99)

	require.NoError(t, err)
	assert.Equal(t, "aprovado", status)
}

func TestCheckoutService_GetOrderStatus_NotFound(t *testing.T) {
	orderRepo := &mocks.MockOrderRepository{
		FindByAppmaxOrderIDAndInstallationFunc: func(_ context.Context, _ int, _ int64) (*models.Order, error) {
			return nil, nil
		},
	}

	svc := mustCheckoutService(t, &mockAppmaxService{}, orderRepo)
	_, err := svc.GetOrderStatus(context.Background(), checkoutTestInst, 99)

	require.ErrorIs(t, err, services.ErrOrderNotFound)
}

func TestCheckoutService_GetOrderStatus_RepoError(t *testing.T) {
	repoErr := errors.New("db connection lost")
	orderRepo := &mocks.MockOrderRepository{
		FindByAppmaxOrderIDAndInstallationFunc: func(_ context.Context, _ int, _ int64) (*models.Order, error) {
			return nil, repoErr
		},
	}

	svc := mustCheckoutService(t, &mockAppmaxService{}, orderRepo)
	_, err := svc.GetOrderStatus(context.Background(), checkoutTestInst, 99)

	require.Error(t, err)
	assert.ErrorContains(t, err, "db connection lost")
}

func TestCheckoutService_ProcessCreditCard_PersistErrorDoesNotFailRequest(t *testing.T) {
	mock := &mockAppmaxService{
		createOrUpdateCustomerFunc: func(_ context.Context, _ *models.Installation, _ services.CustomerInput) (int, error) {
			return 42, nil
		},
		createOrderFunc: func(_ context.Context, _ *models.Installation, _ services.OrderInput) (int, error) {
			return 99, nil
		},
		creditCardFunc: func(_ context.Context, _ *models.Installation, _ services.CreditCardInput) (services.CreditCardResult, error) {
			return services.CreditCardResult{PaymentID: 1, Status: "autorizado"}, nil
		},
	}
	orderRepo := &mocks.MockOrderRepository{
		CreateFunc: func(_ context.Context, _ *models.Order) error {
			return errors.New("db down")
		},
	}

	svc := mustCheckoutService(t, mock, orderRepo)
	result, err := svc.ProcessCreditCard(context.Background(), checkoutTestInst, testCheckoutInput())

	require.NoError(t, err, "persist failure must not bubble up to the caller")
	assert.Equal(t, 99, result.OrderID)
	assert.Equal(t, "autorizado", result.Status)
}

func TestCheckoutServiceConstructor_RejectsNilDependency(t *testing.T) {
	svc, err := services.NewCheckoutService(nil, noopOrderRepo())

	require.Error(t, err)
	assert.Nil(t, svc)
	assert.ErrorIs(t, err, services.ErrNilDependency)
}

package services_test

import (
	"context"
	"errors"
	"testing"

	gatewayappmax "github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/gateway/appmax"
	gatewaycontracts "github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/gateway/appmax/contracts"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/models"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tokenManagerStub struct {
	appToken      string
	merchantToken string
	appErr        error
	merchantErr   error
}

func (m *tokenManagerStub) AppToken(context.Context) (string, error) {
	if m.appErr != nil {
		return "", m.appErr
	}

	return m.appToken, nil
}

func (m *tokenManagerStub) MerchantToken(context.Context, *models.Installation) (string, error) {
	if m.merchantErr != nil {
		return "", m.merchantErr
	}

	return m.merchantToken, nil
}

type appmaxGatewayMock struct {
	authorizeFunc      func(ctx context.Context, appToken, appID, externalKey, callbackURL string) (string, error)
	generateCredsFunc  func(ctx context.Context, appToken, hash string) (string, string, error)
	createCustomerFunc func(ctx context.Context, merchantToken string, req gatewayappmax.CreateCustomerRequest) (int, error)
	createOrderFunc    func(ctx context.Context, merchantToken string, req gatewayappmax.CreateOrderRequest) (int, error)
	getOrderFunc       func(ctx context.Context, merchantToken string, orderID int) (gatewayappmax.GetOrderResponse, error)
	creditCardFunc     func(ctx context.Context, merchantToken string, req gatewayappmax.CreditCardRequest) (gatewayappmax.CreditCardResponse, error)
	pixFunc            func(ctx context.Context, merchantToken string, req gatewayappmax.PixRequest) (gatewayappmax.PixResponse, error)
	boletoFunc         func(ctx context.Context, merchantToken string, req gatewayappmax.BoletoRequest) (gatewayappmax.BoletoResponse, error)
	installmentsFunc   func(ctx context.Context, merchantToken string, req gatewayappmax.InstallmentsRequest) ([]gatewayappmax.InstallmentItem, error)
	refundFunc         func(ctx context.Context, merchantToken string, req gatewayappmax.RefundRequest) error
	addTrackingFunc    func(ctx context.Context, merchantToken string, req gatewayappmax.TrackingRequest) error
	createUpsellFunc   func(ctx context.Context, merchantToken string, req gatewayappmax.UpsellRequest) (gatewayappmax.UpsellResponse, error)
	tokenizeFunc       func(ctx context.Context, merchantToken string, req gatewayappmax.TokenizeRequest) (gatewayappmax.TokenizeResponse, error)
	getTokenFunc       func(ctx context.Context, clientID, clientSecret string) (gatewayappmax.TokenResponse, error)
}

func (m *appmaxGatewayMock) GetToken(ctx context.Context, clientID, clientSecret string) (gatewayappmax.TokenResponse, error) {
	if m.getTokenFunc == nil {
		return gatewayappmax.TokenResponse{}, nil
	}
	return m.getTokenFunc(ctx, clientID, clientSecret)
}

func (m *appmaxGatewayMock) Authorize(ctx context.Context, appToken, appID, externalKey, callbackURL string) (string, error) {
	if m.authorizeFunc == nil {
		return "", nil
	}
	return m.authorizeFunc(ctx, appToken, appID, externalKey, callbackURL)
}

func (m *appmaxGatewayMock) GenerateMerchantCreds(ctx context.Context, appToken, hash string) (string, string, error) {
	if m.generateCredsFunc == nil {
		return "", "", nil
	}
	return m.generateCredsFunc(ctx, appToken, hash)
}

func (m *appmaxGatewayMock) CreateOrUpdateCustomer(ctx context.Context, merchantToken string, req gatewayappmax.CreateCustomerRequest) (int, error) {
	if m.createCustomerFunc == nil {
		return 0, nil
	}
	return m.createCustomerFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) CreateOrder(ctx context.Context, merchantToken string, req gatewayappmax.CreateOrderRequest) (int, error) {
	if m.createOrderFunc == nil {
		return 0, nil
	}
	return m.createOrderFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) GetOrder(ctx context.Context, merchantToken string, orderID int) (gatewayappmax.GetOrderResponse, error) {
	if m.getOrderFunc == nil {
		return gatewayappmax.GetOrderResponse{}, nil
	}
	return m.getOrderFunc(ctx, merchantToken, orderID)
}

func (m *appmaxGatewayMock) CreditCard(ctx context.Context, merchantToken string, req gatewayappmax.CreditCardRequest) (gatewayappmax.CreditCardResponse, error) {
	if m.creditCardFunc == nil {
		return gatewayappmax.CreditCardResponse{}, nil
	}
	return m.creditCardFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) Pix(ctx context.Context, merchantToken string, req gatewayappmax.PixRequest) (gatewayappmax.PixResponse, error) {
	if m.pixFunc == nil {
		return gatewayappmax.PixResponse{}, nil
	}
	return m.pixFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) Boleto(ctx context.Context, merchantToken string, req gatewayappmax.BoletoRequest) (gatewayappmax.BoletoResponse, error) {
	if m.boletoFunc == nil {
		return gatewayappmax.BoletoResponse{}, nil
	}
	return m.boletoFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) Installments(ctx context.Context, merchantToken string, req gatewayappmax.InstallmentsRequest) ([]gatewayappmax.InstallmentItem, error) {
	if m.installmentsFunc == nil {
		return nil, nil
	}
	return m.installmentsFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) Refund(ctx context.Context, merchantToken string, req gatewayappmax.RefundRequest) error {
	if m.refundFunc == nil {
		return nil
	}
	return m.refundFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) AddTracking(ctx context.Context, merchantToken string, req gatewayappmax.TrackingRequest) error {
	if m.addTrackingFunc == nil {
		return nil
	}
	return m.addTrackingFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) CreateUpsell(ctx context.Context, merchantToken string, req gatewayappmax.UpsellRequest) (gatewayappmax.UpsellResponse, error) {
	if m.createUpsellFunc == nil {
		return gatewayappmax.UpsellResponse{}, nil
	}
	return m.createUpsellFunc(ctx, merchantToken, req)
}

func (m *appmaxGatewayMock) Tokenize(ctx context.Context, merchantToken string, req gatewayappmax.TokenizeRequest) (gatewayappmax.TokenizeResponse, error) {
	if m.tokenizeFunc == nil {
		return gatewayappmax.TokenizeResponse{}, nil
	}
	return m.tokenizeFunc(ctx, merchantToken, req)
}

func mustAppmaxService(t *testing.T, tokenManager services.TokenManager, gateway gatewaycontracts.Gateway) services.AppmaxService {
	t.Helper()

	svc, err := services.NewAppmaxServiceWithGateway(tokenManager, gateway)
	require.NoError(t, err)
	return svc
}

func TestAppmaxService_Authorize(t *testing.T) {
	mock := &appmaxGatewayMock{
		authorizeFunc: func(_ context.Context, appToken, appID, externalKey, callbackURL string) (string, error) {
			assert.Equal(t, "app-token", appToken)
			assert.Equal(t, "app-id", appID)
			assert.Equal(t, "external-key", externalKey)
			assert.Equal(t, "https://cb", callbackURL)
			return "hash-1", nil
		},
	}
	svc := mustAppmaxService(t, &tokenManagerStub{appToken: "app-token", merchantToken: "merchant-token"}, mock)

	hash, err := svc.Authorize(context.Background(), "app-id", "external-key", "https://cb")

	require.NoError(t, err)
	assert.Equal(t, "hash-1", hash)
}

func TestAppmaxService_CreateCustomerAndOrderMapDTO(t *testing.T) {
	inst := &models.Installation{ID: 1}

	mock := &appmaxGatewayMock{
		createCustomerFunc: func(_ context.Context, merchantToken string, req gatewayappmax.CreateCustomerRequest) (int, error) {
			assert.Equal(t, "merchant-token", merchantToken)
			assert.Equal(t, "John", req.FirstName)
			require.NotNil(t, req.Address)
			assert.Equal(t, "Porto Alegre", req.Address.City)
			require.NotEmpty(t, req.Products)
			assert.Equal(t, "sku-1", req.Products[0].SKU)
			require.NotNil(t, req.Tracking)
			assert.Equal(t, "google", req.Tracking.UTMSource)
			return 55, nil
		},
		createOrderFunc: func(_ context.Context, merchantToken string, req gatewayappmax.CreateOrderRequest) (int, error) {
			assert.Equal(t, "merchant-token", merchantToken)
			assert.Equal(t, 55, req.CustomerID)
			assert.Equal(t, 10000, req.ProductsValue)
			require.NotEmpty(t, req.Products)
			assert.Equal(t, "sku-1", req.Products[0].SKU)
			return 77, nil
		},
	}
	svc := mustAppmaxService(t, &tokenManagerStub{appToken: "app-token", merchantToken: "merchant-token"}, mock)

	customerID, err := svc.CreateOrUpdateCustomer(context.Background(), inst, services.CustomerInput{
		FirstName:      "John",
		LastName:       "Doe",
		Email:          "john@example.com",
		Phone:          "51999999999",
		DocumentNumber: "123",
		Address:        &services.Address{City: "Porto Alegre"},
		Products:       []services.Product{{SKU: "sku-1", Name: "Product", Quantity: 1, UnitValue: 10000, Type: "digital"}},
		Tracking:       &services.Tracking{UTMSource: "google"},
	})
	require.NoError(t, err)
	assert.Equal(t, 55, customerID)

	orderID, err := svc.CreateOrder(context.Background(), inst, services.OrderInput{
		CustomerID:    customerID,
		ProductsValue: 10000,
		Products:      []services.Product{{SKU: "sku-1", Name: "Product", Quantity: 1, UnitValue: 10000, Type: "digital"}},
	})
	require.NoError(t, err)
	assert.Equal(t, 77, orderID)
}

func TestAppmaxService_PaymentsAndOperations(t *testing.T) {
	inst := &models.Installation{ID: 2}

	mock := &appmaxGatewayMock{
		creditCardFunc: func(_ context.Context, _ string, req gatewayappmax.CreditCardRequest) (gatewayappmax.CreditCardResponse, error) {
			assert.Equal(t, 88, req.OrderID)
			assert.Equal(t, 22, req.CustomerID)
			assert.Equal(t, "4000000000000010", req.PaymentData.CreditCard.Number)
			var out gatewayappmax.CreditCardResponse
			out.Data.Payment.ID = 900
			out.Data.Payment.PayReference = "mock-ref"
			out.Data.Payment.UpsellHash = "up-1"
			return out, nil
		},
		pixFunc: func(_ context.Context, _ string, req gatewayappmax.PixRequest) (gatewayappmax.PixResponse, error) {
			assert.Equal(t, 88, req.OrderID)
			assert.Equal(t, "123", req.PaymentData.Pix.DocumentNumber)
			var out gatewayappmax.PixResponse
			out.Data.Payment.QRCode = "qr"
			out.Data.Payment.EMV = "emv"
			return out, nil
		},
		boletoFunc: func(_ context.Context, _ string, req gatewayappmax.BoletoRequest) (gatewayappmax.BoletoResponse, error) {
			assert.Equal(t, 88, req.OrderID)
			var out gatewayappmax.BoletoResponse
			out.Data.Payment.PDFURL = "https://pdf"
			out.Data.Payment.Digitavel = "digitavel"
			return out, nil
		},
		installmentsFunc: func(_ context.Context, _ string, req gatewayappmax.InstallmentsRequest) ([]gatewayappmax.InstallmentItem, error) {
			assert.True(t, req.Settings)
			return []gatewayappmax.InstallmentItem{
				{Installments: 1, Value: 100, TotalValue: 100},
				{Installments: 2, Value: 55, TotalValue: 110},
			}, nil
		},
		refundFunc: func(_ context.Context, _ string, req gatewayappmax.RefundRequest) error {
			assert.Equal(t, "total", req.Type)
			return nil
		},
		getOrderFunc: func(_ context.Context, _ string, orderID int) (gatewayappmax.GetOrderResponse, error) {
			assert.Equal(t, 88, orderID)
			var out gatewayappmax.GetOrderResponse
			out.Data.Order.ID = 88
			out.Data.Order.Status = "paid"
			out.Data.Order.TotalPaid = 10000
			out.Data.Customer.Name = "John Doe"
			return out, nil
		},
		addTrackingFunc: func(_ context.Context, _ string, req gatewayappmax.TrackingRequest) error {
			assert.Equal(t, "BR123", req.ShippingTrackingCode)
			return nil
		},
		createUpsellFunc: func(_ context.Context, _ string, req gatewayappmax.UpsellRequest) (gatewayappmax.UpsellResponse, error) {
			require.NotEmpty(t, req.Products)
			var out gatewayappmax.UpsellResponse
			out.Data.Order.ID = 99
			out.Data.Order.Status = "approved"
			return out, nil
		},
		tokenizeFunc: func(_ context.Context, _ string, req gatewayappmax.TokenizeRequest) (gatewayappmax.TokenizeResponse, error) {
			assert.Equal(t, "4000000000000010", req.PaymentData.CreditCard.Number)
			var out gatewayappmax.TokenizeResponse
			out.Data.Token = "tok-1"
			return out, nil
		},
	}
	svc := mustAppmaxService(t, &tokenManagerStub{appToken: "app-token", merchantToken: "merchant-token"}, mock)

	cc, err := svc.CreditCard(context.Background(), inst, services.CreditCardInput{
		OrderID:         88,
		CustomerID:      22,
		Number:          "4000000000000010",
		Installments:    1,
		ExpirationMonth: "12",
		ExpirationYear:  "28",
	})
	require.NoError(t, err)
	assert.Equal(t, 900, cc.PaymentID)
	assert.Equal(t, "aprovado", cc.Status)
	assert.Equal(t, "up-1", cc.UpsellHash)

	pix, err := svc.Pix(context.Background(), inst, services.PixInput{OrderID: 88, DocumentNumber: "123"})
	require.NoError(t, err)
	assert.Equal(t, "qr", pix.QRCode)
	assert.Equal(t, "emv", pix.EMV)

	boleto, err := svc.Boleto(context.Background(), inst, services.BoletoInput{OrderID: 88, DocumentNumber: "123"})
	require.NoError(t, err)
	assert.Equal(t, "https://pdf", boleto.PDFURL)
	assert.Equal(t, "digitavel", boleto.Digitavel)

	installments, err := svc.Installments(context.Background(), inst, services.InstallmentsInput{Installments: 2, TotalValue: 100})
	require.NoError(t, err)
	require.Len(t, installments, 2)
	assert.Equal(t, 2, installments[1].Installments)

	err = svc.Refund(context.Background(), inst, services.RefundInput{OrderID: 88, Type: "total"})
	require.NoError(t, err)

	order, err := svc.GetOrder(context.Background(), inst, 88)
	require.NoError(t, err)
	assert.Equal(t, "John Doe", order.CustomerName)

	err = svc.AddTracking(context.Background(), inst, services.TrackingInput{OrderID: 88, ShippingTrackingCode: "BR123"})
	require.NoError(t, err)

	upsell, err := svc.Upsell(context.Background(), inst, services.UpsellInput{
		UpsellHash:    "up-1",
		ProductsValue: 1000,
		Products:      []services.Product{{SKU: "sku-2", Name: "Upsell", Quantity: 1, UnitValue: 1000, Type: "digital"}},
	})
	require.NoError(t, err)
	assert.Equal(t, 99, upsell.OrderID)

	token, err := svc.Tokenize(context.Background(), inst, services.TokenizeInput{
		Number:          "4000000000000010",
		CVV:             "123",
		ExpirationMonth: "12",
		ExpirationYear:  "28",
		HolderName:      "John",
	})
	require.NoError(t, err)
	assert.Equal(t, "tok-1", token)
}

func TestAppmaxService_MapsDeclinedError(t *testing.T) {
	inst := &models.Installation{ID: 3}

	mock := &appmaxGatewayMock{
		creditCardFunc: func(context.Context, string, gatewayappmax.CreditCardRequest) (gatewayappmax.CreditCardResponse, error) {
			return gatewayappmax.CreditCardResponse{}, errors.New("status 402")
		},
	}
	svc := mustAppmaxService(t, &tokenManagerStub{appToken: "app-token", merchantToken: "merchant-token"}, mock)

	_, err := svc.CreditCard(context.Background(), inst, services.CreditCardInput{OrderID: 1, CustomerID: 1})

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrPaymentDeclined)
}

func TestAppmaxService_PropagatesTokenManagerError(t *testing.T) {
	svc := mustAppmaxService(t, &tokenManagerStub{appErr: errors.New("auth unavailable"), merchantErr: errors.New("merchant unavailable")}, &appmaxGatewayMock{})

	_, err := svc.Authorize(context.Background(), "app-id", "ext", "cb")
	require.Error(t, err)
	assert.ErrorContains(t, err, "auth unavailable")

	_, err = svc.CreateOrder(context.Background(), &models.Installation{ID: 1}, services.OrderInput{})
	require.Error(t, err)
	assert.ErrorContains(t, err, "merchant unavailable")
}

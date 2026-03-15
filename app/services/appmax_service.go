package services

import (
	"context"
	"fmt"
	"strings"

	gatewaycontracts "github.com/geovanne-gallinati/AppStoreAppDemo/app/gateway/appmax/contracts"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
)

type Address struct {
	Postcode   string
	Street     string
	Number     string
	Complement string
	District   string
	City       string
	State      string
}

type Product struct {
	SKU       string
	Name      string
	Quantity  int
	UnitValue int
	Type      string
}

type Tracking struct {
	UTMSource   string
	UTMCampaign string
}

type CustomerInput struct {
	FirstName      string
	LastName       string
	Email          string
	Phone          string
	DocumentNumber string
	Address        *Address
	IP             string
	Products       []Product
	Tracking       *Tracking
}

type OrderInput struct {
	CustomerID    int
	ProductsValue int
	DiscountValue int
	ShippingValue int
	Products      []Product
}

type CreditCardInput struct {
	OrderID              int
	CustomerID           int
	Token                string
	UpsellHash           string
	Number               string
	CVV                  string
	ExpirationMonth      string
	ExpirationYear       string
	HolderDocumentNumber string
	HolderName           string
	Installments         int
	SoftDescriptor       string
}

type CreditCardResult struct {
	PaymentID  int
	Status     string
	UpsellHash string
}

type PixInput struct {
	OrderID        int
	DocumentNumber string
}

type PixResult struct {
	QRCode string
	EMV    string
}

type BoletoInput struct {
	OrderID        int
	DocumentNumber string
}

type BoletoResult struct {
	PDFURL    string
	Digitavel string
}

type InstallmentsInput struct {
	Installments int
	TotalValue   int
}

type RefundInput struct {
	OrderID int
	Type    string
	Value   int
}

type TrackingInput struct {
	OrderID              int
	ShippingTrackingCode string
}

type UpsellInput struct {
	UpsellHash    string
	ProductsValue int
	Products      []Product
}

type UpsellResult struct {
	OrderID int
	Status  string
}

type GetOrderResult struct {
	OrderID      int
	Status       string
	TotalPaid    int
	CustomerName string
}

type TokenizeInput struct {
	Number          string
	CVV             string
	ExpirationMonth string
	ExpirationYear  string
	HolderName      string
}

type AppmaxInstallmentItem struct {
	Installments int
	Value        float64
	TotalValue   float64
}

type AppmaxService interface {
	Authorize(ctx context.Context, appID, externalKey, callbackURL string) (string, error)
	GenerateMerchantCreds(ctx context.Context, hash string) (string, string, error)
	VerifyMerchantCreds(ctx context.Context, clientID, clientSecret string) error
	CreateOrUpdateCustomer(ctx context.Context, inst *models.Installation, input CustomerInput) (int, error)
	CreateOrder(ctx context.Context, inst *models.Installation, input OrderInput) (int, error)
	CreditCard(ctx context.Context, inst *models.Installation, input CreditCardInput) (CreditCardResult, error)
	Pix(ctx context.Context, inst *models.Installation, input PixInput) (PixResult, error)
	Boleto(ctx context.Context, inst *models.Installation, input BoletoInput) (BoletoResult, error)
	Installments(ctx context.Context, inst *models.Installation, input InstallmentsInput) ([]AppmaxInstallmentItem, error)
	Refund(ctx context.Context, inst *models.Installation, input RefundInput) error
	GetOrder(ctx context.Context, inst *models.Installation, orderID int) (GetOrderResult, error)
	AddTracking(ctx context.Context, inst *models.Installation, input TrackingInput) error
	Upsell(ctx context.Context, inst *models.Installation, input UpsellInput) (UpsellResult, error)
	Tokenize(ctx context.Context, inst *models.Installation, input TokenizeInput) (string, error)
}

type appmaxService struct {
	tokenManager TokenManager
	gateway      gatewaycontracts.Gateway
}

var _ AppmaxService = (*appmaxService)(nil)

func NewAppmaxServiceWithGateway(tokenManager TokenManager, gateway gatewaycontracts.Gateway) (AppmaxService, error) {
	if tokenManager == nil || gateway == nil {
		return nil, fmt.Errorf("new appmax service: %w", ErrNilDependency)
	}

	return &appmaxService{
		tokenManager: tokenManager,
		gateway:      gateway,
	}, nil
}

func (s *appmaxService) withAppToken(ctx context.Context, fn func(appToken string) error) error {
	token, err := s.tokenManager.AppToken(ctx)
	if err != nil {
		return err
	}
	return fn(token)
}

func (s *appmaxService) withMerchantToken(ctx context.Context, inst *models.Installation, fn func(merchantToken string) error) error {
	token, err := s.tokenManager.MerchantToken(ctx, inst)
	if err != nil {
		return err
	}
	return fn(token)
}

func (s *appmaxService) Authorize(ctx context.Context, appID, externalKey, callbackURL string) (string, error) {
	var hash string
	err := s.withAppToken(ctx, func(appToken string) error {
		var callErr error
		hash, callErr = s.gateway.Authorize(ctx, appToken, appID, externalKey, callbackURL)
		return callErr
	})
	if err != nil {
		return "", fmt.Errorf("appmax authorize: %w", err)
	}
	return hash, nil
}

func (s *appmaxService) VerifyMerchantCreds(ctx context.Context, clientID, clientSecret string) error {
	_, err := s.gateway.GetToken(ctx, clientID, clientSecret)
	if err != nil {
		return fmt.Errorf("appmax verify merchant creds: %w", err)
	}
	return nil
}

func (s *appmaxService) GenerateMerchantCreds(ctx context.Context, hash string) (string, string, error) {
	var clientID, clientSecret string
	err := s.withAppToken(ctx, func(appToken string) error {
		var callErr error
		clientID, clientSecret, callErr = s.gateway.GenerateMerchantCreds(ctx, appToken, hash)
		return callErr
	})
	if err != nil {
		return "", "", fmt.Errorf("appmax generate merchant creds: %w", err)
	}
	return clientID, clientSecret, nil
}

func (s *appmaxService) CreateOrUpdateCustomer(ctx context.Context, inst *models.Installation, input CustomerInput) (int, error) {
	var customerID int
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		req := gatewaycontracts.CreateCustomerRequest{
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			Email:          input.Email,
			Phone:          input.Phone,
			DocumentNumber: input.DocumentNumber,
			IP:             input.IP,
			Products:       toGatewayProducts(input.Products),
		}
		req.Address = toGatewayAddress(input.Address)
		req.Tracking = toGatewayTracking(input.Tracking)
		var callErr error
		customerID, callErr = s.gateway.CreateOrUpdateCustomer(ctx, merchantToken, req)
		return callErr
	})
	if err != nil {
		return 0, fmt.Errorf("appmax create customer: %w", err)
	}
	return customerID, nil
}

func (s *appmaxService) CreateOrder(ctx context.Context, inst *models.Installation, input OrderInput) (int, error) {
	var orderID int
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		req := gatewaycontracts.CreateOrderRequest{
			CustomerID:    input.CustomerID,
			ProductsValue: input.ProductsValue,
			DiscountValue: input.DiscountValue,
			ShippingValue: input.ShippingValue,
			Products:      toGatewayProducts(input.Products),
		}
		var callErr error
		orderID, callErr = s.gateway.CreateOrder(ctx, merchantToken, req)
		return callErr
	})
	if err != nil {
		return 0, fmt.Errorf("appmax create order: %w", err)
	}
	return orderID, nil
}

func (s *appmaxService) CreditCard(ctx context.Context, inst *models.Installation, input CreditCardInput) (CreditCardResult, error) {
	var result CreditCardResult
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		req := gatewaycontracts.CreditCardRequest{
			OrderID:    input.OrderID,
			CustomerID: input.CustomerID,
		}
		req.PaymentData.CreditCard = gatewaycontracts.CreditCardData{
			Token:                input.Token,
			UpsellHash:           input.UpsellHash,
			Number:               input.Number,
			CVV:                  input.CVV,
			ExpirationMonth:      input.ExpirationMonth,
			ExpirationYear:       input.ExpirationYear,
			HolderDocumentNumber: input.HolderDocumentNumber,
			HolderName:           input.HolderName,
			Installments:         input.Installments,
			SoftDescriptor:       input.SoftDescriptor,
		}
		resp, callErr := s.gateway.CreditCard(ctx, merchantToken, req)
		if callErr != nil {
			if isDeclinedError(callErr) {
				return ErrPaymentDeclined
			}
			return callErr
		}
		result = CreditCardResult{
			PaymentID:  resp.Data.Payment.ID,
			Status:     resp.Data.Payment.Status,
			UpsellHash: resp.Data.Payment.UpsellHash,
		}
		return nil
	})
	if err != nil {
		return CreditCardResult{}, fmt.Errorf("appmax credit card: %w", err)
	}
	return result, nil
}

func (s *appmaxService) Pix(ctx context.Context, inst *models.Installation, input PixInput) (PixResult, error) {
	var result PixResult
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		req := gatewaycontracts.PixRequest{OrderID: input.OrderID}
		req.PaymentData.Pix.DocumentNumber = input.DocumentNumber
		resp, callErr := s.gateway.Pix(ctx, merchantToken, req)
		if callErr != nil {
			return callErr
		}
		result = PixResult{
			QRCode: resp.Data.Payment.QRCode,
			EMV:    resp.Data.Payment.EMV,
		}
		return nil
	})
	if err != nil {
		return PixResult{}, fmt.Errorf("appmax pix: %w", err)
	}
	return result, nil
}

func (s *appmaxService) Boleto(ctx context.Context, inst *models.Installation, input BoletoInput) (BoletoResult, error) {
	var result BoletoResult
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		req := gatewaycontracts.BoletoRequest{OrderID: input.OrderID}
		req.PaymentData.Boleto.DocumentNumber = input.DocumentNumber
		resp, callErr := s.gateway.Boleto(ctx, merchantToken, req)
		if callErr != nil {
			return callErr
		}
		result = BoletoResult{
			PDFURL:    resp.Data.Payment.PDFURL,
			Digitavel: resp.Data.Payment.Digitavel,
		}
		return nil
	})
	if err != nil {
		return BoletoResult{}, fmt.Errorf("appmax boleto: %w", err)
	}
	return result, nil
}

func (s *appmaxService) Installments(ctx context.Context, inst *models.Installation, input InstallmentsInput) ([]AppmaxInstallmentItem, error) {
	var items []gatewaycontracts.InstallmentItem
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		req := gatewaycontracts.InstallmentsRequest{
			Installments: input.Installments,
			TotalValue:   input.TotalValue,
			Settings:     true,
		}
		var callErr error
		items, callErr = s.gateway.Installments(ctx, merchantToken, req)
		return callErr
	})
	if err != nil {
		return nil, fmt.Errorf("appmax installments: %w", err)
	}
	return fromGatewayInstallmentItems(items), nil
}

func (s *appmaxService) Refund(ctx context.Context, inst *models.Installation, input RefundInput) error {
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		return s.gateway.Refund(ctx, merchantToken, gatewaycontracts.RefundRequest{
			OrderID: input.OrderID,
			Type:    input.Type,
			Value:   input.Value,
		})
	})
	if err != nil {
		return fmt.Errorf("appmax refund: %w", err)
	}
	return nil
}

func (s *appmaxService) GetOrder(ctx context.Context, inst *models.Installation, orderID int) (GetOrderResult, error) {
	var result GetOrderResult
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		resp, callErr := s.gateway.GetOrder(ctx, merchantToken, orderID)
		if callErr != nil {
			return callErr
		}
		result = GetOrderResult{
			OrderID:      resp.Data.Order.ID,
			Status:       resp.Data.Order.Status,
			TotalPaid:    resp.Data.Order.TotalPaid,
			CustomerName: resp.Data.Customer.Name,
		}
		return nil
	})
	if err != nil {
		return GetOrderResult{}, fmt.Errorf("appmax get order: %w", err)
	}
	return result, nil
}

func (s *appmaxService) AddTracking(ctx context.Context, inst *models.Installation, input TrackingInput) error {
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		return s.gateway.AddTracking(ctx, merchantToken, gatewaycontracts.TrackingRequest{
			OrderID:              input.OrderID,
			ShippingTrackingCode: input.ShippingTrackingCode,
		})
	})
	if err != nil {
		return fmt.Errorf("appmax add tracking: %w", err)
	}
	return nil
}

func (s *appmaxService) Upsell(ctx context.Context, inst *models.Installation, input UpsellInput) (UpsellResult, error) {
	var result UpsellResult
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		resp, callErr := s.gateway.CreateUpsell(ctx, merchantToken, gatewaycontracts.UpsellRequest{
			UpsellHash:    input.UpsellHash,
			ProductsValue: input.ProductsValue,
			Products:      toGatewayProducts(input.Products),
		})
		if callErr != nil {
			return callErr
		}
		result = UpsellResult{
			OrderID: resp.Data.Order.ID,
			Status:  resp.Data.Order.Status,
		}
		return nil
	})
	if err != nil {
		return UpsellResult{}, fmt.Errorf("appmax upsell: %w", err)
	}
	return result, nil
}

func (s *appmaxService) Tokenize(ctx context.Context, inst *models.Installation, input TokenizeInput) (string, error) {
	var token string
	err := s.withMerchantToken(ctx, inst, func(merchantToken string) error {
		var req gatewaycontracts.TokenizeRequest
		req.PaymentData.CreditCard.Number = input.Number
		req.PaymentData.CreditCard.CVV = input.CVV
		req.PaymentData.CreditCard.ExpirationMonth = input.ExpirationMonth
		req.PaymentData.CreditCard.ExpirationYear = input.ExpirationYear
		req.PaymentData.CreditCard.HolderName = input.HolderName
		resp, callErr := s.gateway.Tokenize(ctx, merchantToken, req)
		if callErr != nil {
			return callErr
		}
		token = resp.Data.Token
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("appmax tokenize: %w", err)
	}
	return token, nil
}

func toGatewayAddress(input *Address) *gatewaycontracts.Address {
	if input == nil {
		return nil
	}

	return &gatewaycontracts.Address{
		Postcode:   input.Postcode,
		Street:     input.Street,
		Number:     input.Number,
		Complement: input.Complement,
		District:   input.District,
		City:       input.City,
		State:      input.State,
	}
}

func toGatewayProducts(input []Product) []gatewaycontracts.Product {
	if len(input) == 0 {
		return nil
	}

	out := make([]gatewaycontracts.Product, len(input))
	for i, product := range input {
		out[i] = gatewaycontracts.Product{
			SKU:       product.SKU,
			Name:      product.Name,
			Quantity:  product.Quantity,
			UnitValue: product.UnitValue,
			Type:      product.Type,
		}
	}

	return out
}

func toGatewayTracking(input *Tracking) *gatewaycontracts.Tracking {
	if input == nil {
		return nil
	}

	return &gatewaycontracts.Tracking{
		UTMSource:   input.UTMSource,
		UTMCampaign: input.UTMCampaign,
	}
}

func fromGatewayInstallmentItems(input []gatewaycontracts.InstallmentItem) []AppmaxInstallmentItem {
	if len(input) == 0 {
		return nil
	}

	out := make([]AppmaxInstallmentItem, len(input))
	for i, item := range input {
		out[i] = AppmaxInstallmentItem{
			Installments: item.Installments,
			Value:        item.Value,
			TotalValue:   item.TotalValue,
		}
	}

	return out
}

func isDeclinedError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not authorized") ||
		strings.Contains(msg, "declined") ||
		strings.Contains(msg, "recusado") ||
		strings.Contains(msg, "status 402") ||
		strings.Contains(msg, "status 422")
}

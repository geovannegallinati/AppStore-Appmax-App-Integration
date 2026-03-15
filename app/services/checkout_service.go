package services

import (
	"context"
	"fmt"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
)

type CheckoutCreditCardInput struct {
	Customer CustomerInput
	Order    OrderInput
	Payment  CreditCardInput
}

type CheckoutPixInput struct {
	Customer       CustomerInput
	Order          OrderInput
	DocumentNumber string
}

type CheckoutBoletoInput struct {
	Customer       CustomerInput
	Order          OrderInput
	DocumentNumber string
}

type CheckoutCreditCardResult struct {
	OrderID    int
	Status     string
	UpsellHash string
}

type CheckoutPixResult struct {
	OrderID int
	QRCode  string
	EMV     string
}

type CheckoutBoletoResult struct {
	OrderID   int
	PDFURL    string
	Digitavel string
}

type CheckoutService interface {
	ProcessCreditCard(ctx context.Context, inst *models.Installation, input CheckoutCreditCardInput) (CheckoutCreditCardResult, error)
	ProcessPix(ctx context.Context, inst *models.Installation, input CheckoutPixInput) (CheckoutPixResult, error)
	ProcessBoleto(ctx context.Context, inst *models.Installation, input CheckoutBoletoInput) (CheckoutBoletoResult, error)
	GetOrderStatus(ctx context.Context, inst *models.Installation, appmaxOrderID int) (string, error)
	GetInstallments(ctx context.Context, inst *models.Installation, input InstallmentsInput) ([]AppmaxInstallmentItem, error)
	ProcessRefund(ctx context.Context, inst *models.Installation, input RefundInput) error
}

type checkoutService struct {
	appmaxSvc AppmaxService
	orderRepo contracts.OrderRepository
	logger    Logger
}

var _ CheckoutService = (*checkoutService)(nil)

func NewCheckoutService(appmaxSvc AppmaxService, orderRepo contracts.OrderRepository, logger ...Logger) (CheckoutService, error) {
	if appmaxSvc == nil || orderRepo == nil {
		return nil, fmt.Errorf("new checkout service: %w", ErrNilDependency)
	}

	var selectedLogger Logger = noopLogger{}
	if len(logger) > 0 && logger[0] != nil {
		selectedLogger = logger[0]
	}

	return &checkoutService{
		appmaxSvc: appmaxSvc,
		orderRepo: orderRepo,
		logger:    selectedLogger,
	}, nil
}

func (s *checkoutService) ProcessCreditCard(ctx context.Context, inst *models.Installation, input CheckoutCreditCardInput) (CheckoutCreditCardResult, error) {
	customerID, appmaxOrderID, err := s.createCustomerAndOrder(ctx, inst, input.Customer, input.Order, "credit card")
	if err != nil {
		return CheckoutCreditCardResult{}, err
	}

	input.Payment.OrderID = appmaxOrderID
	input.Payment.CustomerID = customerID

	payResult, payErr := s.appmaxSvc.CreditCard(ctx, inst, input.Payment)
	status := "cancelado"
	if payErr == nil {
		status = payResult.Status
	}

	order := &models.Order{
		InstallationID:   inst.ID,
		AppmaxCustomerID: customerID,
		AppmaxOrderID:    appmaxOrderID,
		Status:           status,
		PaymentMethod:    "credit_card",
		UpsellHash:       payResult.UpsellHash,
	}
	s.persistOrderBestEffort(ctx, appmaxOrderID, order)

	if payErr != nil {
		return CheckoutCreditCardResult{}, fmt.Errorf("checkout credit card: %w", payErr)
	}

	return CheckoutCreditCardResult{
		OrderID:    appmaxOrderID,
		Status:     payResult.Status,
		UpsellHash: payResult.UpsellHash,
	}, nil
}

func (s *checkoutService) ProcessPix(ctx context.Context, inst *models.Installation, input CheckoutPixInput) (CheckoutPixResult, error) {
	customerID, appmaxOrderID, err := s.createCustomerAndOrder(ctx, inst, input.Customer, input.Order, "pix")
	if err != nil {
		return CheckoutPixResult{}, err
	}

	pixResult, pixErr := s.appmaxSvc.Pix(ctx, inst, PixInput{
		OrderID:        appmaxOrderID,
		DocumentNumber: input.DocumentNumber,
	})

	status := "pendente"
	if pixErr != nil {
		status = "cancelado"
	}

	order := &models.Order{
		InstallationID:   inst.ID,
		AppmaxCustomerID: customerID,
		AppmaxOrderID:    appmaxOrderID,
		Status:           status,
		PaymentMethod:    "pix",
		PixQRCode:        pixResult.QRCode,
		PixEMV:           pixResult.EMV,
	}
	s.persistOrderBestEffort(ctx, appmaxOrderID, order)

	if pixErr != nil {
		return CheckoutPixResult{}, fmt.Errorf("checkout pix: %w", pixErr)
	}

	return CheckoutPixResult{
		OrderID: appmaxOrderID,
		QRCode:  pixResult.QRCode,
		EMV:     pixResult.EMV,
	}, nil
}

func (s *checkoutService) ProcessBoleto(ctx context.Context, inst *models.Installation, input CheckoutBoletoInput) (CheckoutBoletoResult, error) {
	customerID, appmaxOrderID, err := s.createCustomerAndOrder(ctx, inst, input.Customer, input.Order, "boleto")
	if err != nil {
		return CheckoutBoletoResult{}, err
	}

	boletoResult, boletoErr := s.appmaxSvc.Boleto(ctx, inst, BoletoInput{
		OrderID:        appmaxOrderID,
		DocumentNumber: input.DocumentNumber,
	})

	status := "pendente"
	if boletoErr != nil {
		status = "cancelado"
	}

	order := &models.Order{
		InstallationID:   inst.ID,
		AppmaxCustomerID: customerID,
		AppmaxOrderID:    appmaxOrderID,
		Status:           status,
		PaymentMethod:    "boleto",
		BoletoPDFURL:     boletoResult.PDFURL,
		BoletoDigitavel:  boletoResult.Digitavel,
	}
	s.persistOrderBestEffort(ctx, appmaxOrderID, order)

	if boletoErr != nil {
		return CheckoutBoletoResult{}, fmt.Errorf("checkout boleto: %w", boletoErr)
	}

	return CheckoutBoletoResult{
		OrderID:   appmaxOrderID,
		PDFURL:    boletoResult.PDFURL,
		Digitavel: boletoResult.Digitavel,
	}, nil
}

func (s *checkoutService) GetOrderStatus(ctx context.Context, inst *models.Installation, appmaxOrderID int) (string, error) {
	order, err := s.orderRepo.FindByAppmaxOrderIDAndInstallation(ctx, appmaxOrderID, inst.ID)
	if err != nil {
		return "", fmt.Errorf("get order status: %w", err)
	}
	if order == nil {
		return "", ErrOrderNotFound
	}
	return order.Status, nil
}

func (s *checkoutService) GetInstallments(ctx context.Context, inst *models.Installation, input InstallmentsInput) ([]AppmaxInstallmentItem, error) {
	return s.appmaxSvc.Installments(ctx, inst, input)
}

func (s *checkoutService) ProcessRefund(ctx context.Context, inst *models.Installation, input RefundInput) error {
	return s.appmaxSvc.Refund(ctx, inst, input)
}

func (s *checkoutService) createCustomerAndOrder(ctx context.Context, inst *models.Installation, customer CustomerInput, order OrderInput, flow string) (int, int, error) {
	customerID, err := s.appmaxSvc.CreateOrUpdateCustomer(ctx, inst, customer)
	if err != nil {
		return 0, 0, fmt.Errorf("checkout %s: %w", flow, err)
	}

	order.CustomerID = customerID
	appmaxOrderID, err := s.appmaxSvc.CreateOrder(ctx, inst, order)
	if err != nil {
		return 0, 0, fmt.Errorf("checkout %s: %w", flow, err)
	}

	return customerID, appmaxOrderID, nil
}

func (s *checkoutService) persistOrderBestEffort(ctx context.Context, appmaxOrderID int, order *models.Order) {
	if createErr := s.orderRepo.Create(ctx, order); createErr != nil {
		s.logger.Warningf("checkout_service: failed to persist order %d: %v", appmaxOrderID, createErr)
	}
}

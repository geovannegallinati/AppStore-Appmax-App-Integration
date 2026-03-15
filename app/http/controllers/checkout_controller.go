package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/http/middleware"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/http/requests"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/http/responses"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

type CheckoutController struct {
	checkoutSvc services.CheckoutService
}

func NewCheckoutController(checkoutSvc services.CheckoutService) (*CheckoutController, error) {
	if checkoutSvc == nil {
		return nil, fmt.Errorf("new checkout controller: %w", ErrNilDependency)
	}

	return &CheckoutController{checkoutSvc: checkoutSvc}, nil
}

func installationFromCtx(ctx http.Context) (*models.Installation, bool) {
	inst, ok := ctx.Value(middleware.InstallationContextKey).(*models.Installation)
	return inst, ok && inst != nil
}

func (c *CheckoutController) PayCreditCard(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	var body requests.CheckoutCreditCardRequest
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	input := services.CheckoutCreditCardInput{
		Customer: toCustomerInput(body.Customer),
		Order:    toOrderInput(body.Order),
		Payment: services.CreditCardInput{
			Token:                body.Payment.Token,
			UpsellHash:           body.Payment.UpsellHash,
			Number:               body.Payment.Number,
			CVV:                  body.Payment.CVV,
			ExpirationMonth:      body.Payment.ExpirationMonth,
			ExpirationYear:       body.Payment.ExpirationYear,
			HolderDocumentNumber: body.Payment.HolderDocumentNumber,
			HolderName:           body.Payment.HolderName,
			Installments:         body.Payment.Installments,
			SoftDescriptor:       body.Payment.SoftDescriptor,
		},
	}

	result, err := c.checkoutSvc.ProcessCreditCard(ctx.Context(), inst, input)
	if err != nil {
		if errors.Is(err, services.ErrPaymentDeclined) {
			return ctx.Response().Json(422, responses.MessageResponse{Message: "payment declined"})
		}
		facades.Log().Errorf("checkout_controller: credit card failed: %v", err)
		return ctx.Response().Json(502, responses.MessageResponse{Message: "payment processing failed"})
	}

	return ctx.Response().Json(200, responses.CheckoutCreditCardResponse{
		OrderID: result.OrderID,
		Status:  result.Status,
	})
}

func (c *CheckoutController) PayPix(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	var body requests.CheckoutPixRequest
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	input := services.CheckoutPixInput{
		Customer:       toCustomerInput(body.Customer),
		Order:          toOrderInput(body.Order),
		DocumentNumber: body.DocumentNumber,
	}

	result, err := c.checkoutSvc.ProcessPix(ctx.Context(), inst, input)
	if err != nil {
		facades.Log().Errorf("checkout_controller: pix failed: %v", err)
		return ctx.Response().Json(502, responses.MessageResponse{Message: "payment processing failed"})
	}

	return ctx.Response().Json(200, responses.CheckoutPixResponse{
		OrderID: result.OrderID,
		QRCode:  result.QRCode,
		EMV:     result.EMV,
	})
}

func (c *CheckoutController) PayBoleto(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	var body requests.CheckoutBoletoRequest
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	input := services.CheckoutBoletoInput{
		Customer:       toCustomerInput(body.Customer),
		Order:          toOrderInput(body.Order),
		DocumentNumber: body.DocumentNumber,
	}

	result, err := c.checkoutSvc.ProcessBoleto(ctx.Context(), inst, input)
	if err != nil {
		facades.Log().Errorf("checkout_controller: boleto failed: %v", err)
		return ctx.Response().Json(502, responses.MessageResponse{Message: "payment processing failed"})
	}

	return ctx.Response().Json(200, responses.CheckoutBoletoResponse{
		OrderID:   result.OrderID,
		PDFURL:    result.PDFURL,
		Digitavel: result.Digitavel,
	})
}

func (c *CheckoutController) Status(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	orderIDStr := ctx.Request().Route("order_id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid order_id"})
	}

	status, err := c.checkoutSvc.GetOrderStatus(ctx.Context(), inst, orderID)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			return ctx.Response().Json(404, responses.MessageResponse{Message: "order not found"})
		}
		facades.Log().Errorf("checkout_controller: get status failed: %v", err)
		return ctx.Response().Json(500, responses.MessageResponse{Message: "internal server error"})
	}

	return ctx.Response().Json(200, responses.StatusResponse{Status: status})
}

func (c *CheckoutController) Installments(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	totalValueStr := ctx.Request().Query("total_value", "0")
	installmentsStr := ctx.Request().Query("installments", "12")

	totalValue, err := strconv.Atoi(totalValueStr)
	if err != nil || totalValue <= 0 {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "total_value is required"})
	}

	installments, err := strconv.Atoi(installmentsStr)
	if err != nil {
		installments = 12
	}
	if installments <= 0 {
		installments = 12
	}

	items, err := c.checkoutSvc.GetInstallments(ctx.Context(), inst, services.InstallmentsInput{
		Installments: installments,
		TotalValue:   totalValue,
	})
	if err != nil {
		facades.Log().Errorf("checkout_controller: installments failed: %v", err)
		return ctx.Response().Json(502, responses.MessageResponse{Message: "failed to fetch installments"})
	}

	return ctx.Response().Json(200, items)
}

func (c *CheckoutController) Refund(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	var body requests.CheckoutRefundRequest
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	if body.OrderID <= 0 {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "order_id is required"})
	}

	err := c.checkoutSvc.ProcessRefund(ctx.Context(), inst, services.RefundInput{
		OrderID: body.OrderID,
		Type:    body.Type,
		Value:   body.Value,
	})
	if err != nil {
		facades.Log().Errorf("checkout_controller: refund failed: %v", err)
		return ctx.Response().Json(502, responses.MessageResponse{Message: "refund request failed"})
	}

	return ctx.Response().Json(200, responses.MessageResponse{Message: "Refund request accepted"})
}

func toServiceAddress(address *requests.Address) *services.Address {
	if address == nil {
		return nil
	}

	return &services.Address{
		Postcode:   address.Postcode,
		Street:     address.Street,
		Number:     address.Number,
		Complement: address.Complement,
		District:   address.District,
		City:       address.City,
		State:      address.State,
	}
}

func toServiceProducts(products []requests.Product) []services.Product {
	out := make([]services.Product, len(products))
	for i, product := range products {
		out[i] = services.Product{
			SKU:       product.SKU,
			Name:      product.Name,
			Quantity:  product.Quantity,
			UnitValue: product.UnitValue,
			Type:      product.Type,
		}
	}

	return out
}

func toCustomerInput(customer requests.Customer) services.CustomerInput {
	return services.CustomerInput{
		FirstName:      customer.FirstName,
		LastName:       customer.LastName,
		Email:          customer.Email,
		Phone:          customer.Phone,
		DocumentNumber: customer.DocumentNumber,
		IP:             customer.IP,
		Address:        toServiceAddress(customer.Address),
	}
}

func toOrderInput(order requests.Order) services.OrderInput {
	return services.OrderInput{
		ProductsValue: order.ProductsValue,
		DiscountValue: order.DiscountValue,
		ShippingValue: order.ShippingValue,
		Products:      toServiceProducts(order.Products),
	}
}

package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/middleware"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/requests"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/responses"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/models"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
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

func (c *CheckoutController) CreateOrder(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	var body requests.CheckoutCreateOrderRequest
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	result, err := c.checkoutSvc.CreateCustomerAndOrder(ctx.Context(), inst, toCustomerInput(body.Customer), toOrderInput(body.Order))
	if err != nil {
		facades.Log().Errorf("checkout_controller: create order failed: %v", err)
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "order creation failed")})
	}

	return ctx.Response().Json(200, responses.CheckoutCreateOrderResponse{
		CustomerID: result.CustomerID,
		OrderID:    result.OrderID,
	})
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
		CustomerID: derefInt(body.CustomerID),
		OrderID:    derefInt(body.OrderID),
		Customer:   toCustomerInput(body.Customer),
		Order:      toOrderInput(body.Order),
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
		Subscription: toServiceSubscription(body.Subscription),
	}

	result, err := c.checkoutSvc.ProcessCreditCard(ctx.Context(), inst, input)
	if err != nil {
		if errors.Is(err, services.ErrPaymentDeclined) {
			return ctx.Response().Json(422, responses.MessageResponse{Message: "payment declined"})
		}
		facades.Log().Errorf("checkout_controller: credit card failed: %v", err)
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "payment processing failed")})
	}

	return ctx.Response().Json(200, responses.CheckoutCreditCardResponse{
		OrderID:    result.OrderID,
		Status:     result.Status,
		UpsellHash: result.UpsellHash,
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
		CustomerID:     derefInt(body.CustomerID),
		OrderID:        derefInt(body.OrderID),
		Customer:       toCustomerInput(body.Customer),
		Order:          toOrderInput(body.Order),
		DocumentNumber: body.DocumentNumber,
		Subscription:   toServiceSubscription(body.Subscription),
	}

	result, err := c.checkoutSvc.ProcessPix(ctx.Context(), inst, input)
	if err != nil {
		facades.Log().Errorf("checkout_controller: pix failed: %v", err)
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "payment processing failed")})
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
		CustomerID:     derefInt(body.CustomerID),
		OrderID:        derefInt(body.OrderID),
		Customer:       toCustomerInput(body.Customer),
		Order:          toOrderInput(body.Order),
		DocumentNumber: body.DocumentNumber,
	}

	result, err := c.checkoutSvc.ProcessBoleto(ctx.Context(), inst, input)
	if err != nil {
		facades.Log().Errorf("checkout_controller: boleto failed: %v", err)
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "payment processing failed")})
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
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "failed to fetch installments")})
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
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: RefundErrorMessage(err)})
	}

	return ctx.Response().Json(200, responses.MessageResponse{Message: "Refund request accepted"})
}

func (c *CheckoutController) Tokenize(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	var body requests.CheckoutTokenizeRequest
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	token, err := c.checkoutSvc.Tokenize(ctx.Context(), inst, services.TokenizeInput{
		Number:          body.Number,
		CVV:             body.CVV,
		ExpirationMonth: body.ExpirationMonth,
		ExpirationYear:  body.ExpirationYear,
		HolderName:      body.HolderName,
	})
	if err != nil {
		facades.Log().Errorf("checkout_controller: tokenize failed: %v", err)
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "tokenization failed")})
	}

	return ctx.Response().Json(200, responses.CheckoutTokenizeResponse{Token: token})
}

func (c *CheckoutController) AddTracking(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	var body requests.CheckoutTrackingRequest
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	if body.OrderID <= 0 {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "order_id is required"})
	}

	err := c.checkoutSvc.AddTracking(ctx.Context(), inst, services.TrackingInput{
		OrderID:              body.OrderID,
		ShippingTrackingCode: body.ShippingTrackingCode,
	})
	if err != nil {
		facades.Log().Errorf("checkout_controller: tracking failed: %v", err)
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "tracking update failed")})
	}

	return ctx.Response().Json(200, responses.CheckoutTrackingResponse{Message: "tracking accepted"})
}

func (c *CheckoutController) Upsell(ctx http.Context) http.Response {
	inst, ok := installationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	var body requests.CheckoutUpsellRequest
	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	result, err := c.checkoutSvc.ProcessUpsell(ctx.Context(), inst, services.UpsellInput{
		UpsellHash:    body.UpsellHash,
		ProductsValue: body.ProductsValue,
		Products:      toServiceProducts(body.Products),
	})
	if err != nil {
		facades.Log().Errorf("checkout_controller: upsell failed: %v", err)
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "upsell failed")})
	}

	return ctx.Response().Json(200, responses.CheckoutUpsellResponse{
		Message:     result.Message,
		RedirectURL: result.RedirectURL,
	})
}

func toServiceSubscription(sub *requests.CheckoutSubscription) *services.Subscription {
	if sub == nil {
		return nil
	}
	return &services.Subscription{
		Interval:      sub.Interval,
		IntervalCount: sub.IntervalCount,
	}
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

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func toOrderInput(order requests.Order) services.OrderInput {
	return services.OrderInput{
		ProductsValue: order.ProductsValue,
		DiscountValue: order.DiscountValue,
		ShippingValue: order.ShippingValue,
		Products:      toServiceProducts(order.Products),
	}
}

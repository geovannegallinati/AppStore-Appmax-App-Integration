package appmax

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) CreditCard(ctx context.Context, merchantToken string, req CreditCardRequest) (CreditCardResponse, error) {
	out, err := doAndDecode[CreditCardResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/payments/credit-card", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return CreditCardResponse{}, fmt.Errorf("credit card payment: %w", err)
	}

	return out, nil
}

func (c *Client) Pix(ctx context.Context, merchantToken string, req PixRequest) (PixResponse, error) {
	out, err := doAndDecode[PixResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/payments/pix", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return PixResponse{}, fmt.Errorf("pix payment: %w", err)
	}

	return out, nil
}

func (c *Client) Boleto(ctx context.Context, merchantToken string, req BoletoRequest) (BoletoResponse, error) {
	out, err := doAndDecode[BoletoResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/payments/boleto", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return BoletoResponse{}, fmt.Errorf("boleto payment: %w", err)
	}

	return out, nil
}

func (c *Client) Installments(ctx context.Context, merchantToken string, req InstallmentsRequest) ([]InstallmentItem, error) {
	out, err := doAndDecode[InstallmentsResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/payments/installments", req, merchantToken, http.StatusOK)
	if err != nil {
		return nil, fmt.Errorf("installments: %w", err)
	}

	return out.Data, nil
}

func (c *Client) Refund(ctx context.Context, merchantToken string, req RefundRequest) error {
	_, err := doAndDecode[struct{}](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/orders/refund-request", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("refund: %w", err)
	}

	return nil
}

func (c *Client) AddTracking(ctx context.Context, merchantToken string, req TrackingRequest) error {
	_, err := doAndDecode[struct{}](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/orders/shipping-tracking-code", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return fmt.Errorf("add tracking: %w", err)
	}

	return nil
}

func (c *Client) CreateUpsell(ctx context.Context, merchantToken string, req UpsellRequest) (UpsellResponse, error) {
	out, err := doAndDecode[UpsellResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/orders/upsell", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return UpsellResponse{}, fmt.Errorf("upsell: %w", err)
	}

	return out, nil
}

func (c *Client) Tokenize(ctx context.Context, merchantToken string, req TokenizeRequest) (TokenizeResponse, error) {
	out, err := doAndDecode[TokenizeResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/payments/tokenize", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return TokenizeResponse{}, fmt.Errorf("tokenize: %w", err)
	}

	return out, nil
}

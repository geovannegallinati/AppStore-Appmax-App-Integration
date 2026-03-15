package appmax

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) CreateOrUpdateCustomer(ctx context.Context, merchantToken string, req CreateCustomerRequest) (int, error) {
	out, err := doAndDecode[CreateCustomerResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/customers", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return 0, fmt.Errorf("create or update customer: %w", err)
	}

	return out.Data.Customer.ID, nil
}

func (c *Client) CreateOrder(ctx context.Context, merchantToken string, req CreateOrderRequest) (int, error) {
	out, err := doAndDecode[CreateOrderResponse](c, ctx, http.MethodPost, c.apiBaseURL+"/v1/orders", req, merchantToken, http.StatusOK, http.StatusCreated)
	if err != nil {
		return 0, fmt.Errorf("create order: %w", err)
	}

	return out.Data.Order.ID, nil
}

func (c *Client) GetOrder(ctx context.Context, merchantToken string, orderID int) (GetOrderResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/orders/%d", c.apiBaseURL, orderID)
	out, err := doAndDecode[GetOrderResponse](c, ctx, http.MethodGet, endpoint, nil, merchantToken, http.StatusOK)
	if err != nil {
		return GetOrderResponse{}, fmt.Errorf("get order: %w", err)
	}

	return out, nil
}

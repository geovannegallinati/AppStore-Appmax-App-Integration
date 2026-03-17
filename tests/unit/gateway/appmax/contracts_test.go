package appmax_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	appmax "github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/gateway/appmax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetToken(t *testing.T) {
	client := clientWithTransport(t, func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "https://auth.example.com/oauth2/token", req.URL.String())
		assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))

		return jsonResponse(http.StatusOK, `{"access_token":"tok","token_type":"Bearer","expires_in":3600}`), nil
	})

	resp, err := client.GetToken(context.Background(), "cid", "csecret")

	require.NoError(t, err)
	assert.Equal(t, "tok", resp.AccessToken)
	assert.Equal(t, 3600, resp.ExpiresIn)
}

func TestClient_AuthEndpoints(t *testing.T) {
	client := clientWithTransport(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/app/authorize":
			assert.Equal(t, "Bearer app-token", req.Header.Get("Authorization"))
			return jsonResponse(http.StatusCreated, `{"data":{"token":"hash-123"}}`), nil
		case "/app/client/generate":
			assert.Equal(t, "Bearer app-token", req.Header.Get("Authorization"))
			return jsonResponse(http.StatusOK, `{"data":{"client":{"client_id":"mid","client_secret":"msecret"}}}`), nil
		default:
			t.Fatalf("unexpected path: %s", req.URL.Path)
			return nil, nil
		}
	})

	hash, err := client.Authorize(context.Background(), "app-token", "app-1", "ext-1", "https://callback")
	require.NoError(t, err)
	assert.Equal(t, "hash-123", hash)

	clientID, clientSecret, err := client.GenerateMerchantCreds(context.Background(), "app-token", "hash-123")
	require.NoError(t, err)
	assert.Equal(t, "mid", clientID)
	assert.Equal(t, "msecret", clientSecret)
}

func TestClient_CustomerOrderEndpoints(t *testing.T) {
	client := clientWithTransport(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/customers":
			assert.Equal(t, http.MethodPost, req.Method)
			return jsonResponse(http.StatusCreated, `{"data":{"customer":{"id":55}}}`), nil
		case "/v1/orders":
			assert.Equal(t, http.MethodPost, req.Method)
			return jsonResponse(http.StatusOK, `{"data":{"order":{"id":88,"status":"pending"}}}`), nil
		case "/v1/orders/88":
			assert.Equal(t, http.MethodGet, req.Method)
			return jsonResponse(http.StatusOK, `{"data":{"order":{"id":88,"status":"paid","total_paid":10000},"customer":{"name":"John"}}}`), nil
		default:
			t.Fatalf("unexpected path: %s", req.URL.Path)
			return nil, nil
		}
	})

	customerID, err := client.CreateOrUpdateCustomer(context.Background(), "merchant-token", appmax.CreateCustomerRequest{FirstName: "John"})
	require.NoError(t, err)
	assert.Equal(t, 55, customerID)

	orderID, err := client.CreateOrder(context.Background(), "merchant-token", appmax.CreateOrderRequest{CustomerID: customerID})
	require.NoError(t, err)
	assert.Equal(t, 88, orderID)

	order, err := client.GetOrder(context.Background(), "merchant-token", orderID)
	require.NoError(t, err)
	assert.Equal(t, "paid", order.Data.Order.Status)
	assert.Equal(t, "John", order.Data.Customer.Name)
}

func TestClient_PaymentEndpoints(t *testing.T) {
	client := clientWithTransport(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/payments/credit-card":
			return jsonResponse(http.StatusCreated, `{"data":{"payment":{"id":1,"status":"approved","upsell_hash":"up"}}}`), nil
		case "/v1/payments/pix":
			return jsonResponse(http.StatusOK, `{"data":{"payment":{"pix_qrcode":"qr","pix_emv":"emv"}}}`), nil
		case "/v1/payments/boleto":
			return jsonResponse(http.StatusCreated, `{"data":{"payment":{"boleto_link_pdf":"https://pdf","boleto_digitable_line":"123"}}}`), nil
		case "/v1/payments/installments":
			return jsonResponse(http.StatusOK, `{"data":{"parcels":{"1":100.00}}}`), nil
		case "/v1/orders/refund-request":
			return jsonResponse(http.StatusCreated, `{}`), nil
		case "/v1/orders/shipping-tracking-code":
			return jsonResponse(http.StatusOK, `{}`), nil
		case "/v1/orders/upsell":
			return jsonResponse(http.StatusOK, `{"data":{"order":{"id":99,"status":"approved"}}}`), nil
		case "/v1/payments/tokenize":
			return jsonResponse(http.StatusCreated, `{"data":{"token":"tok-1"}}`), nil
		default:
			t.Fatalf("unexpected path: %s", req.URL.Path)
			return nil, nil
		}
	})

	cc, err := client.CreditCard(context.Background(), "merchant-token", appmax.CreditCardRequest{})
	require.NoError(t, err)
	assert.Equal(t, "approved", cc.Data.Payment.Status)

	pix, err := client.Pix(context.Background(), "merchant-token", appmax.PixRequest{})
	require.NoError(t, err)
	assert.Equal(t, "qr", pix.Data.Payment.QRCode)

	boleto, err := client.Boleto(context.Background(), "merchant-token", appmax.BoletoRequest{})
	require.NoError(t, err)
	assert.Equal(t, "123", boleto.Data.Payment.Digitavel)

	installments, err := client.Installments(context.Background(), "merchant-token", appmax.InstallmentsRequest{})
	require.NoError(t, err)
	require.Len(t, installments, 1)
	assert.Equal(t, 1, installments[0].Installments)

	require.NoError(t, client.Refund(context.Background(), "merchant-token", appmax.RefundRequest{OrderID: 1, Type: "total"}))
	require.NoError(t, client.AddTracking(context.Background(), "merchant-token", appmax.TrackingRequest{OrderID: 1, ShippingTrackingCode: "BR123"}))

	upsell, err := client.CreateUpsell(context.Background(), "merchant-token", appmax.UpsellRequest{})
	require.NoError(t, err)
	assert.Equal(t, 99, upsell.Data.Order.ID)

	tokenize, err := client.Tokenize(context.Background(), "merchant-token", appmax.TokenizeRequest{})
	require.NoError(t, err)
	assert.Equal(t, "tok-1", tokenize.Data.Token)
}

func TestClient_DoRetriesAndContextCancel(t *testing.T) {
	attempts := 0
	client, err := appmax.NewClientWithOptions("https://auth.example.com", "https://api.example.com", appmax.ClientOptions{
		HTTPClient: &http.Client{
			Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
				attempts++
				if attempts == 1 {
					return nil, errors.New("temporary network error")
				}
				return jsonResponse(http.StatusOK, `{"data":{"order":{"id":1,"status":"ok"}}}`), nil
			}),
		},
		RetryMax:  1,
		RetryWait: 1 * time.Millisecond,
	})
	require.NoError(t, err)

	_, err = client.GetOrder(context.Background(), "merchant-token", 1)
	require.NoError(t, err)
	assert.Equal(t, 2, attempts)

	cancelClient, err := appmax.NewClientWithOptions("https://auth.example.com", "https://api.example.com", appmax.ClientOptions{
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network down")
		}),
		RetryMax:  2,
		RetryWait: 50 * time.Millisecond,
	})
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = cancelClient.GetOrder(ctx, "merchant-token", 1)
	require.Error(t, err)
	assert.ErrorContains(t, err, "context canceled")
}

func TestClient_DecodeAndStatusErrors(t *testing.T) {
	t.Run("status error", func(t *testing.T) {
		client := clientWithTransport(t, func(*http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusBadRequest, `{"message":"invalid"}`), nil
		})

		_, err := client.CreditCard(context.Background(), "merchant-token", appmax.CreditCardRequest{})
		require.Error(t, err)
		assert.ErrorContains(t, err, "unexpected status 400")
	})

	t.Run("decode error", func(t *testing.T) {
		client := clientWithTransport(t, func(*http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `not-json`), nil
		})

		_, err := client.Pix(context.Background(), "merchant-token", appmax.PixRequest{})
		require.Error(t, err)
		assert.ErrorContains(t, err, "decode response")
	})
}

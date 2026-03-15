package contracts

import "context"

type TokenGateway interface {
	GetToken(ctx context.Context, clientID, clientSecret string) (TokenResponse, error)
}

type Gateway interface {
	TokenGateway
	Authorize(ctx context.Context, appToken, appID, externalKey, callbackURL string) (string, error)
	GenerateMerchantCreds(ctx context.Context, appToken, hash string) (string, string, error)
	CreateOrUpdateCustomer(ctx context.Context, merchantToken string, req CreateCustomerRequest) (int, error)
	CreateOrder(ctx context.Context, merchantToken string, req CreateOrderRequest) (int, error)
	GetOrder(ctx context.Context, merchantToken string, orderID int) (GetOrderResponse, error)
	CreditCard(ctx context.Context, merchantToken string, req CreditCardRequest) (CreditCardResponse, error)
	Pix(ctx context.Context, merchantToken string, req PixRequest) (PixResponse, error)
	Boleto(ctx context.Context, merchantToken string, req BoletoRequest) (BoletoResponse, error)
	Installments(ctx context.Context, merchantToken string, req InstallmentsRequest) ([]InstallmentItem, error)
	Refund(ctx context.Context, merchantToken string, req RefundRequest) error
	AddTracking(ctx context.Context, merchantToken string, req TrackingRequest) error
	CreateUpsell(ctx context.Context, merchantToken string, req UpsellRequest) (UpsellResponse, error)
	Tokenize(ctx context.Context, merchantToken string, req TokenizeRequest) (TokenizeResponse, error)
}

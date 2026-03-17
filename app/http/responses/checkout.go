package responses

type CheckoutCreateOrderResponse struct {
	CustomerID int `json:"customer_id"`
	OrderID    int `json:"order_id"`
}

type CheckoutCreditCardResponse struct {
	OrderID    int    `json:"order_id"`
	Status     string `json:"status"`
	UpsellHash string `json:"upsell_hash,omitempty"`
}

type CheckoutPixResponse struct {
	OrderID int    `json:"order_id"`
	QRCode  string `json:"qr_code"`
	EMV     string `json:"emv"`
}

type CheckoutBoletoResponse struct {
	OrderID   int    `json:"order_id"`
	PDFURL    string `json:"pdf_url"`
	Digitavel string `json:"digitavel"`
}

type CheckoutTokenizeResponse struct {
	Token string `json:"token"`
}

type CheckoutTrackingResponse struct {
	Message string `json:"message"`
}

type CheckoutUpsellResponse struct {
	Message     string `json:"message"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

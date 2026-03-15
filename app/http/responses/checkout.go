package responses

type CheckoutCreditCardResponse struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
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

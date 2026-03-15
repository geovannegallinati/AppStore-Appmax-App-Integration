package requests

type Address struct {
	Postcode   string `json:"postcode"`
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state"`
}

type Product struct {
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitValue int    `json:"unit_value"`
	Type      string `json:"type"`
}

type Customer struct {
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	DocumentNumber string   `json:"document_number"`
	IP             string   `json:"ip"`
	Address        *Address `json:"address"`
}

type Order struct {
	ProductsValue int       `json:"products_value"`
	DiscountValue int       `json:"discount_value"`
	ShippingValue int       `json:"shipping_value"`
	Products      []Product `json:"products"`
}

type CheckoutCreditCardPayment struct {
	Token                string `json:"token"`
	UpsellHash           string `json:"upsell_hash"`
	Number               string `json:"number"`
	CVV                  string `json:"cvv"`
	ExpirationMonth      string `json:"expiration_month"`
	ExpirationYear       string `json:"expiration_year"`
	HolderDocumentNumber string `json:"holder_document_number"`
	HolderName           string `json:"holder_name"`
	Installments         int    `json:"installments"`
	SoftDescriptor       string `json:"soft_descriptor"`
}

type CheckoutCreditCardRequest struct {
	Customer Customer                  `json:"customer"`
	Order    Order                     `json:"order"`
	Payment  CheckoutCreditCardPayment `json:"payment"`
}

type CheckoutPixRequest struct {
	Customer       Customer `json:"customer"`
	Order          Order    `json:"order"`
	DocumentNumber string   `json:"document_number"`
}

type CheckoutBoletoRequest struct {
	Customer       Customer `json:"customer"`
	Order          Order    `json:"order"`
	DocumentNumber string   `json:"document_number"`
}

type CheckoutRefundRequest struct {
	OrderID int    `json:"order_id"`
	Type    string `json:"type"`
	Value   int    `json:"value"`
}

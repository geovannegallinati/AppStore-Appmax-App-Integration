package contracts

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type AuthorizeRequest struct {
	AppID       string `json:"app_id"`
	ExternalKey string `json:"external_key"`
	URLCallback string `json:"url_callback"`
}

type AuthorizeResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type GenerateCredsRequest struct {
	Token string `json:"token"`
}

type GenerateCredsResponse struct {
	Data struct {
		Client struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"client"`
	} `json:"data"`
}

type Address struct {
	Postcode   string `json:"postcode,omitempty"`
	Street     string `json:"street,omitempty"`
	Number     string `json:"number,omitempty"`
	Complement string `json:"complement,omitempty"`
	District   string `json:"district,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
}

type Product struct {
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitValue int    `json:"unit_value,omitempty"`
	Type      string `json:"type,omitempty"`
}

type Tracking struct {
	UTMSource   string `json:"utm_source,omitempty"`
	UTMCampaign string `json:"utm_campaign,omitempty"`
}

type CreateCustomerRequest struct {
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	DocumentNumber string    `json:"document_number,omitempty"`
	Address        *Address  `json:"address,omitempty"`
	IP             string    `json:"ip"`
	Products       []Product `json:"products,omitempty"`
	Tracking       *Tracking `json:"tracking,omitempty"`
}

type CreateCustomerResponse struct {
	Data struct {
		Customer struct {
			ID int `json:"id"`
		} `json:"customer"`
	} `json:"data"`
}

type CreateOrderRequest struct {
	CustomerID    int       `json:"customer_id"`
	ProductsValue int       `json:"products_value,omitempty"`
	DiscountValue int       `json:"discount_value"`
	ShippingValue int       `json:"shipping_value"`
	Products      []Product `json:"products"`
}

type CreateOrderResponse struct {
	Data struct {
		Order struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		} `json:"order"`
	} `json:"data"`
}

type GetOrderResponse struct {
	Data struct {
		Order struct {
			ID        int    `json:"id"`
			Status    string `json:"status"`
			TotalPaid int    `json:"total_paid"`
			Amounts   struct {
				SubTotal       int `json:"sub_total"`
				ShippingValue  int `json:"shipping_value"`
				Discount       int `json:"discount"`
				InstallmentFee int `json:"installment_fee"`
			} `json:"amounts"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"order"`
		Customer struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Email          string `json:"email"`
			DocumentNumber string `json:"document_number"`
		} `json:"customer"`
		Payment struct {
			Method             string `json:"method"`
			Installments       int    `json:"installments"`
			InstallmentsAmount int    `json:"installments_amount"`
			Card               struct {
				Brand  string `json:"brand"`
				Number string `json:"number"`
			} `json:"card"`
			PaidAt string `json:"paid_at"`
		} `json:"payment"`
	} `json:"data"`
}

type CreditCardData struct {
	Token                string `json:"token,omitempty"`
	UpsellHash           string `json:"upsell_hash,omitempty"`
	Number               string `json:"number,omitempty"`
	CVV                  string `json:"cvv,omitempty"`
	ExpirationMonth      string `json:"expiration_month,omitempty"`
	ExpirationYear       string `json:"expiration_year,omitempty"`
	HolderDocumentNumber string `json:"holder_document_number,omitempty"`
	HolderName           string `json:"holder_name,omitempty"`
	Installments         int    `json:"installments"`
	SoftDescriptor       string `json:"soft_descriptor,omitempty"`
}

type CreditCardRequest struct {
	OrderID     int `json:"order_id"`
	CustomerID  int `json:"customer_id"`
	PaymentData struct {
		CreditCard CreditCardData `json:"credit_card"`
	} `json:"payment_data"`
}

type CreditCardResponse struct {
	Data struct {
		Payment struct {
			ID         int    `json:"id"`
			Status     string `json:"status"`
			UpsellHash string `json:"upsell_hash,omitempty"`
		} `json:"payment"`
	} `json:"data"`
}

type PixRequest struct {
	OrderID     int `json:"order_id"`
	PaymentData struct {
		Pix struct {
			DocumentNumber string `json:"document_number"`
		} `json:"pix"`
	} `json:"payment_data"`
}

type PixResponse struct {
	Data struct {
		Payment struct {
			QRCode string `json:"qr_code"`
			EMV    string `json:"emv"`
		} `json:"payment"`
	} `json:"data"`
}

type BoletoRequest struct {
	OrderID     int `json:"order_id"`
	PaymentData struct {
		Boleto struct {
			DocumentNumber string `json:"document_number"`
		} `json:"boleto"`
	} `json:"payment_data"`
}

type BoletoResponse struct {
	Data struct {
		Payment struct {
			PDFURL    string `json:"pdf_url"`
			Digitavel string `json:"digitavel"`
		} `json:"payment"`
	} `json:"data"`
}

type InstallmentsRequest struct {
	Installments int  `json:"installments"`
	TotalValue   int  `json:"total_value"`
	Settings     bool `json:"settings"`
}

type InstallmentItem struct {
	Installments int     `json:"installments"`
	Value        float64 `json:"value"`
	TotalValue   float64 `json:"total_value"`
}

type InstallmentsResponse struct {
	Data []InstallmentItem `json:"data"`
}

type RefundRequest struct {
	OrderID int    `json:"order_id"`
	Type    string `json:"type"`
	Value   int    `json:"value,omitempty"`
}

type UpsellRequest struct {
	UpsellHash    string    `json:"upsell_hash"`
	ProductsValue int       `json:"products_value"`
	Products      []Product `json:"products"`
}

type UpsellResponse struct {
	Data struct {
		Order struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		} `json:"order"`
	} `json:"data"`
}

type TrackingRequest struct {
	OrderID              int    `json:"order_id"`
	ShippingTrackingCode string `json:"shipping_tracking_code"`
}

type TokenizeRequest struct {
	PaymentData struct {
		CreditCard struct {
			Number          string `json:"number"`
			CVV             string `json:"cvv"`
			ExpirationMonth string `json:"expiration_month"`
			ExpirationYear  string `json:"expiration_year"`
			HolderName      string `json:"holder_name"`
		} `json:"credit_card"`
	} `json:"payment_data"`
}

type TokenizeResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

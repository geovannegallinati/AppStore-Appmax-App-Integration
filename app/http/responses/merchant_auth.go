package responses

type MerchantTokenSyncResponse struct {
	MerchantBearerToken string `json:"merchant_bearer_token"`
	ExternalKey         string `json:"external_key"`
}

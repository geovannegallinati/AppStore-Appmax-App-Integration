package responses

type MerchantTokenSyncResponse struct {
	ExternalKey         string `json:"external_key"`
	ExternalID          string `json:"external_id"`
	MerchantClientID    string `json:"merchant_client_id"`
	MerchantBearerToken string `json:"merchant_bearer_token"`
	TokenType           string `json:"token_type"`
}

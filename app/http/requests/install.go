package requests

type InstallCallbackRequest struct {
	AppID                string `json:"app_id"`
	ExternalKey          string `json:"external_key"`
	ClientKey            string `json:"client_key"`
	MerchantClientID     string `json:"client_id"`
	MerchantClientSecret string `json:"client_secret"`
}

//go:build appmax_live

package appmax

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const credsFile = ".sandbox-creds.json"

type SandboxCreds struct {
	AppToken             string    `json:"app_token"`
	AppTokenExpiry       time.Time `json:"app_token_expiry"`
	MerchantClientID     string    `json:"merchant_client_id"`
	MerchantClientSecret string    `json:"merchant_client_secret"`
	MerchantToken        string    `json:"merchant_token"`
	MerchantTokenExpiry  time.Time `json:"merchant_token_expiry"`
}

func (c *SandboxCreds) IsAppTokenValid() bool {
	return c.AppToken != "" && time.Now().Before(c.AppTokenExpiry)
}

func (c *SandboxCreds) IsMerchantCredentialsReady() bool {
	return c.MerchantClientID != "" && c.MerchantClientSecret != ""
}

func (c *SandboxCreds) IsMerchantTokenValid() bool {
	return c.MerchantToken != "" && time.Now().Before(c.MerchantTokenExpiry)
}

func credsPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), credsFile)
}

func loadCreds() (*SandboxCreds, error) {
	data, err := os.ReadFile(credsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return &SandboxCreds{}, nil
		}
		return nil, err
	}
	var c SandboxCreds
	if jsonErr := json.Unmarshal(data, &c); jsonErr != nil {
		return &SandboxCreds{}, nil
	}
	return &c, nil
}

func saveCreds(c *SandboxCreds) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credsPath(), data, 0600)
}

package bootstrap

import (
	"fmt"
	"os"
	"strings"
)

const (
	defaultAuthURL  = "https://auth.appmax.com.br"
	defaultAPIURL   = "https://api.appmax.com.br"
	defaultAdminURL = "https://admin.appmax.com.br"
)

type AppmaxConfig struct {
	AuthURL          string
	APIURL           string
	AdminURL         string
	AppPublicURL     string
	AppClientID      string
	AppClientSecret  string
	AppIDUUID        string
	AppIDNumeric     string
}

func LoadAppmaxConfigFromEnv() (AppmaxConfig, error) {
	authURL := os.Getenv("APPMAX_AUTH_URL")
	if strings.TrimSpace(authURL) == "" {
		authURL = defaultAuthURL
	}

	apiURL := os.Getenv("APPMAX_API_URL")
	if strings.TrimSpace(apiURL) == "" {
		apiURL = defaultAPIURL
	}

	adminURL := os.Getenv("APPMAX_ADMIN_URL")
	if strings.TrimSpace(adminURL) == "" {
		adminURL = defaultAdminURL
	}

	appPublicURL := os.Getenv("NGROK_URL")
	if strings.TrimSpace(appPublicURL) == "" {
		appPublicURL = os.Getenv("APP_URL")
	}

	cfg := AppmaxConfig{
		AuthURL:         authURL,
		APIURL:          apiURL,
		AdminURL:        adminURL,
		AppPublicURL:    appPublicURL,
		AppClientID:     os.Getenv("APPMAX_CLIENT_ID"),
		AppClientSecret: os.Getenv("APPMAX_CLIENT_SECRET"),
		AppIDUUID:       os.Getenv("APPMAX_APP_ID_UUID"),
		AppIDNumeric:    os.Getenv("APPMAX_APP_ID_NUMERIC"),
	}

	if err := cfg.Validate(); err != nil {
		return AppmaxConfig{}, err
	}

	return cfg, nil
}

func (c AppmaxConfig) Validate() error {
	var missing []string
	if strings.TrimSpace(c.AppClientID) == "" {
		missing = append(missing, "APPMAX_CLIENT_ID")
	}
	if strings.TrimSpace(c.AppClientSecret) == "" {
		missing = append(missing, "APPMAX_CLIENT_SECRET")
	}
	if strings.TrimSpace(c.AppIDUUID) == "" {
		missing = append(missing, "APPMAX_APP_ID_UUID")
	}
	if strings.TrimSpace(c.AppIDNumeric) == "" {
		missing = append(missing, "APPMAX_APP_ID_NUMERIC")
	}
	if strings.TrimSpace(c.AppPublicURL) == "" {
		missing = append(missing, "APP_URL or NGROK_URL")
	}

	if len(missing) > 0 {
		return fmt.Errorf("invalid appmax configuration, missing env vars: %s", strings.Join(missing, ", "))
	}

	return nil
}

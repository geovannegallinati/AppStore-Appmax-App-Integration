package bootstrap_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/bootstrap"
)

func TestLoadAppmaxConfigFromEnv_MissingVariables(t *testing.T) {
	t.Setenv("APPMAX_CLIENT_ID", "")
	t.Setenv("APPMAX_CLIENT_SECRET", "")

	cfg, err := bootstrap.LoadAppmaxConfigFromEnv()

	require.Error(t, err)
	assert.Equal(t, bootstrap.AppmaxConfig{}, cfg)
	assert.ErrorContains(t, err, "APPMAX_CLIENT_ID")
	assert.ErrorContains(t, err, "APPMAX_CLIENT_SECRET")
}

func TestLoadAppmaxConfigFromEnv_DefaultURLs(t *testing.T) {
	t.Setenv("APPMAX_AUTH_URL", "")
	t.Setenv("APPMAX_API_URL", "")
	t.Setenv("APPMAX_CLIENT_ID", "cid")
	t.Setenv("APPMAX_CLIENT_SECRET", "csecret")
	t.Setenv("APP_ID_UUID", "test-app-uuid")
	t.Setenv("APPMAX_APP_ID_NUMERIC", "123")
	t.Setenv("APP_URL", "https://app.example.com")

	cfg, err := bootstrap.LoadAppmaxConfigFromEnv()

	require.NoError(t, err)
	assert.Equal(t, "https://auth.appmax.com.br", cfg.AuthURL)
	assert.Equal(t, "https://api.appmax.com.br", cfg.APIURL)
}

func TestLoadAppmaxConfigFromEnv_Success(t *testing.T) {
	t.Setenv("APPMAX_AUTH_URL", "https://auth.example.com")
	t.Setenv("APPMAX_API_URL", "https://api.example.com")
	t.Setenv("APPMAX_CLIENT_ID", "cid")
	t.Setenv("APPMAX_CLIENT_SECRET", "csecret")
	t.Setenv("APP_ID_UUID", "test-app-uuid")
	t.Setenv("APPMAX_APP_ID_NUMERIC", "123")
	t.Setenv("APP_URL", "https://app.example.com")

	cfg, err := bootstrap.LoadAppmaxConfigFromEnv()

	require.NoError(t, err)
	assert.Equal(t, "https://auth.example.com", cfg.AuthURL)
	assert.Equal(t, "https://api.example.com", cfg.APIURL)
	assert.Equal(t, "cid", cfg.AppClientID)
	assert.Equal(t, "csecret", cfg.AppClientSecret)
}

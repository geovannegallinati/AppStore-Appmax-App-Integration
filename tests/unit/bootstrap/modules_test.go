package bootstrap_test

import (
	"testing"

	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/bootstrap"
)

type fakeORM struct{}

func (fakeORM) Query() contractsorm.Query {
	return nil
}

func validCfg() bootstrap.AppmaxConfig {
	return bootstrap.AppmaxConfig{
		AuthURL:         "https://auth.example.com",
		APIURL:          "https://api.example.com",
		AdminURL:        "https://admin.example.com",
		AppPublicURL:    "https://app.example.com",
		AppClientID:     "cid",
		AppClientSecret: "csecret",
		AppIDUUID:       "test-app-uuid",
		AppIDNumeric:    "123",
	}
}

func TestNewGatewayModule_InvalidConfig(t *testing.T) {
	module, err := bootstrap.NewGatewayModule(bootstrap.AppmaxConfig{})

	require.Error(t, err)
	assert.Nil(t, module)
}

func TestNewRepositoryModule_NilORM(t *testing.T) {
	module, err := bootstrap.NewRepositoryModule(nil)

	require.Error(t, err)
	assert.Nil(t, module)
	assert.ErrorIs(t, err, bootstrap.ErrNilDependency)
}

func TestNewServiceModule_NilDependencies(t *testing.T) {
	cfg := validCfg()

	module, err := bootstrap.NewServiceModule(cfg, nil, nil)

	require.Error(t, err)
	assert.Nil(t, module)
	assert.ErrorIs(t, err, bootstrap.ErrNilDependency)
}

func TestNewControllerModule_NilDependencies(t *testing.T) {
	cfg := validCfg()

	module, err := bootstrap.NewControllerModule(cfg, nil)

	require.Error(t, err)
	assert.Nil(t, module)
	assert.ErrorIs(t, err, bootstrap.ErrNilDependency)
}

func TestNewHTTPDependenciesWithORM_NilORM(t *testing.T) {
	cfg := validCfg()

	deps, err := bootstrap.NewHTTPDependenciesWithORM(cfg, nil)

	require.Error(t, err)
	assert.Nil(t, deps)
	assert.ErrorIs(t, err, bootstrap.ErrNilDependency)
}

func TestNewGatewayModule_Success(t *testing.T) {
	module, err := bootstrap.NewGatewayModule(validCfg())

	require.NoError(t, err)
	assert.NotNil(t, module)
	assert.NotNil(t, module.AppmaxGateway)
}

func TestNewRepositoryModule_Success(t *testing.T) {
	module, err := bootstrap.NewRepositoryModule(fakeORM{})

	require.NoError(t, err)
	assert.NotNil(t, module)
	assert.NotNil(t, module.InstallationRepository)
	assert.NotNil(t, module.OrderRepository)
	assert.NotNil(t, module.WebhookEventRepository)
}

func TestNewServiceModule_Success(t *testing.T) {
	gateways, err := bootstrap.NewGatewayModule(validCfg())
	require.NoError(t, err)

	repositories, err := bootstrap.NewRepositoryModule(fakeORM{})
	require.NoError(t, err)

	module, err := bootstrap.NewServiceModule(validCfg(), gateways, repositories)

	require.NoError(t, err)
	assert.NotNil(t, module)
	assert.NotNil(t, module.TokenManager)
	assert.NotNil(t, module.AppmaxService)
	assert.NotNil(t, module.InstallService)
	assert.NotNil(t, module.CheckoutService)
	assert.NotNil(t, module.WebhookService)
}

func TestNewControllerModule_Success(t *testing.T) {
	gateways, err := bootstrap.NewGatewayModule(validCfg())
	require.NoError(t, err)

	repositories, err := bootstrap.NewRepositoryModule(fakeORM{})
	require.NoError(t, err)

	services, err := bootstrap.NewServiceModule(validCfg(), gateways, repositories)
	require.NoError(t, err)

	module, err := bootstrap.NewControllerModule(validCfg(), services)

	require.NoError(t, err)
	assert.NotNil(t, module)
	assert.NotNil(t, module.InstallController)
	assert.NotNil(t, module.MerchantAuthController)
	assert.NotNil(t, module.CheckoutController)
	assert.NotNil(t, module.WebhookController)
}

func TestNewHTTPDependenciesWithORM_Success(t *testing.T) {
	deps, err := bootstrap.NewHTTPDependenciesWithORM(validCfg(), fakeORM{})

	require.NoError(t, err)
	assert.NotNil(t, deps)
	assert.NotNil(t, deps.InstallController)
	assert.NotNil(t, deps.MerchantAuthController)
	assert.NotNil(t, deps.CheckoutController)
	assert.NotNil(t, deps.WebhookController)
	assert.NotNil(t, deps.InstallationRepository)
}

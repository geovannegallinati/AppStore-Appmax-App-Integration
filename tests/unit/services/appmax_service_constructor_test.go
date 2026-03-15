package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

type staticTokenManager struct{}

func (m *staticTokenManager) AppToken(_ context.Context) (string, error) {
	return "app-token", nil
}

func (m *staticTokenManager) MerchantToken(_ context.Context, _ *models.Installation) (string, error) {
	return "merchant-token", nil
}

func TestAppmaxServiceConstructor_RejectsNilDependency(t *testing.T) {
	svc, err := services.NewAppmaxServiceWithGateway(nil, nil)

	require.Error(t, err)
	assert.Nil(t, svc)
	assert.ErrorIs(t, err, services.ErrNilDependency)
}

func TestAppmaxServiceConstructor_Success(t *testing.T) {
	svc, err := services.NewAppmaxServiceWithGateway(&staticTokenManager{}, &appmaxGatewayMock{})

	require.NoError(t, err)
	assert.NotNil(t, svc)
}

package repositories_test

import (
	"testing"

	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories"
)

type fakeORM struct{}

func (fakeORM) Query() contractsorm.Query {
	return nil
}

func TestRepositoryConstructors_RejectNilORM(t *testing.T) {
	installationRepo, err := repositories.NewInstallationRepository(nil)
	require.Error(t, err)
	assert.Nil(t, installationRepo)

	orderRepo, err := repositories.NewOrderRepository(nil)
	require.Error(t, err)
	assert.Nil(t, orderRepo)

	eventRepo, err := repositories.NewWebhookEventRepository(nil)
	require.Error(t, err)
	assert.Nil(t, eventRepo)
}

func TestRepositoryConstructors_Success(t *testing.T) {
	orm := fakeORM{}

	installationRepo, err := repositories.NewInstallationRepository(orm)
	require.NoError(t, err)
	assert.NotNil(t, installationRepo)

	orderRepo, err := repositories.NewOrderRepository(orm)
	require.NoError(t, err)
	assert.NotNil(t, orderRepo)

	eventRepo, err := repositories.NewWebhookEventRepository(orm)
	require.NoError(t, err)
	assert.NotNil(t, eventRepo)
}

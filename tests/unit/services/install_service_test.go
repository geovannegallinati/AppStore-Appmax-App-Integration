package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
	"github.com/geovanne-gallinati/AppStoreAppDemo/tests/unit/mocks"
)

var testUpsertInput = services.UpsertInstallationInput{
	AppID:                "app-1",
	ExternalKey:          "store-abc",
	MerchantClientID:     "mc-id",
	MerchantClientSecret: "mc-secret",
}

func TestInstallService_Upsert_CreatesNew(t *testing.T) {
	var created *models.Installation
	repo := &mocks.MockInstallationRepository{
		FindByExternalKeyFunc: func(_ context.Context, _ string) (*models.Installation, error) {
			return nil, nil
		},
		CreateFunc: func(_ context.Context, inst *models.Installation) error {
			inst.ID = 10
			created = inst
			return nil
		},
	}

	svc, err := services.NewInstallService(repo)
	require.NoError(t, err)
	inst, wasCreated, err := svc.Upsert(context.Background(), testUpsertInput)

	require.NoError(t, err)
	assert.True(t, wasCreated)
	assert.Equal(t, int64(10), inst.ID)
	assert.Equal(t, "app-1", inst.AppID)
	assert.Equal(t, "store-abc", inst.ExternalKey)
	assert.Equal(t, "mc-id", inst.MerchantClientID)
	assert.Equal(t, "mc-secret", inst.MerchantClientSecret)
	assert.NotNil(t, created)
}

func TestInstallService_Upsert_UpdatesExisting(t *testing.T) {
	existing := &models.Installation{
		ID:                   5,
		AppID:                "app-1",
		ExternalKey:          "store-abc",
		MerchantClientID:     "old-id",
		MerchantClientSecret: "old-secret",
	}
	var saved *models.Installation
	repo := &mocks.MockInstallationRepository{
		FindByExternalKeyFunc: func(_ context.Context, _ string) (*models.Installation, error) {
			return existing, nil
		},
		SaveFunc: func(_ context.Context, inst *models.Installation) error {
			saved = inst
			return nil
		},
	}

	svc, err := services.NewInstallService(repo)
	require.NoError(t, err)
	inst, wasCreated, err := svc.Upsert(context.Background(), testUpsertInput)

	require.NoError(t, err)
	assert.False(t, wasCreated)
	assert.Equal(t, int64(5), inst.ID)
	assert.Equal(t, "mc-id", inst.MerchantClientID)
	assert.Equal(t, "mc-secret", inst.MerchantClientSecret)
	assert.NotNil(t, saved)
	assert.Equal(t, "mc-id", saved.MerchantClientID)
}

func TestInstallService_Upsert_FindError(t *testing.T) {
	repo := &mocks.MockInstallationRepository{
		FindByExternalKeyFunc: func(_ context.Context, _ string) (*models.Installation, error) {
			return nil, errors.New("db timeout")
		},
	}

	svc, err := services.NewInstallService(repo)
	require.NoError(t, err)
	_, _, err = svc.Upsert(context.Background(), testUpsertInput)

	require.Error(t, err)
	assert.ErrorContains(t, err, "db timeout")
}

func TestInstallService_Upsert_CreateError(t *testing.T) {
	repo := &mocks.MockInstallationRepository{
		FindByExternalKeyFunc: func(_ context.Context, _ string) (*models.Installation, error) {
			return nil, nil
		},
		CreateFunc: func(_ context.Context, _ *models.Installation) error {
			return errors.New("unique constraint violation")
		},
	}

	svc, err := services.NewInstallService(repo)
	require.NoError(t, err)
	_, _, err = svc.Upsert(context.Background(), testUpsertInput)

	require.Error(t, err)
	assert.ErrorContains(t, err, "unique constraint violation")
}

func TestInstallService_Upsert_SaveError(t *testing.T) {
	existing := &models.Installation{ID: 5, ExternalKey: "store-abc"}
	repo := &mocks.MockInstallationRepository{
		FindByExternalKeyFunc: func(_ context.Context, _ string) (*models.Installation, error) {
			return existing, nil
		},
		SaveFunc: func(_ context.Context, _ *models.Installation) error {
			return errors.New("disk full")
		},
	}

	svc, err := services.NewInstallService(repo)
	require.NoError(t, err)
	_, _, err = svc.Upsert(context.Background(), testUpsertInput)

	require.Error(t, err)
	assert.ErrorContains(t, err, "disk full")
}

func TestInstallServiceConstructor_RejectsNilDependency(t *testing.T) {
	svc, err := services.NewInstallService(nil)

	require.Error(t, err)
	assert.Nil(t, svc)
	assert.ErrorIs(t, err, services.ErrNilDependency)
}

func TestInstallServiceConstructor_Success(t *testing.T) {
	repo := &mocks.MockInstallationRepository{
		FindByExternalKeyFunc: func(_ context.Context, _ string) (*models.Installation, error) {
			return nil, nil
		},
		CreateFunc: func(_ context.Context, _ *models.Installation) error { return nil },
		SaveFunc:   func(_ context.Context, _ *models.Installation) error { return nil },
	}

	svc, err := services.NewInstallService(repo)

	require.NoError(t, err)
	assert.NotNil(t, svc)
}

package mocks

import (
	"context"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
)

type MockInstallationRepository struct {
	FindByExternalKeyFunc func(ctx context.Context, key string) (*models.Installation, error)
	CreateFunc            func(ctx context.Context, inst *models.Installation) error
	SaveFunc              func(ctx context.Context, inst *models.Installation) error
}

func (m *MockInstallationRepository) FindByExternalKey(ctx context.Context, key string) (*models.Installation, error) {
	return m.FindByExternalKeyFunc(ctx, key)
}

func (m *MockInstallationRepository) Create(ctx context.Context, inst *models.Installation) error {
	return m.CreateFunc(ctx, inst)
}

func (m *MockInstallationRepository) Save(ctx context.Context, inst *models.Installation) error {
	return m.SaveFunc(ctx, inst)
}

package mocks

import (
	"context"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
)

type MockOrderRepository struct {
	FindByAppmaxOrderIDFunc                func(ctx context.Context, appmaxOrderID int) (*models.Order, error)
	FindByAppmaxOrderIDAndInstallationFunc func(ctx context.Context, appmaxOrderID int, installationID int64) (*models.Order, error)
	CreateFunc                             func(ctx context.Context, order *models.Order) error
	SaveFunc                               func(ctx context.Context, order *models.Order) error
}

func (m *MockOrderRepository) FindByAppmaxOrderID(ctx context.Context, appmaxOrderID int) (*models.Order, error) {
	return m.FindByAppmaxOrderIDFunc(ctx, appmaxOrderID)
}

func (m *MockOrderRepository) FindByAppmaxOrderIDAndInstallation(ctx context.Context, appmaxOrderID int, installationID int64) (*models.Order, error) {
	return m.FindByAppmaxOrderIDAndInstallationFunc(ctx, appmaxOrderID, installationID)
}

func (m *MockOrderRepository) Create(ctx context.Context, order *models.Order) error {
	return m.CreateFunc(ctx, order)
}

func (m *MockOrderRepository) Save(ctx context.Context, order *models.Order) error {
	return m.SaveFunc(ctx, order)
}

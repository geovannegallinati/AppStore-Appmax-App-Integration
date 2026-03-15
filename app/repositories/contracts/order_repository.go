package contracts

import (
	"context"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
)

type OrderRepository interface {
	FindByAppmaxOrderID(ctx context.Context, appmaxOrderID int) (*models.Order, error)
	FindByAppmaxOrderIDAndInstallation(ctx context.Context, appmaxOrderID int, installationID int64) (*models.Order, error)
	Create(ctx context.Context, order *models.Order) error
	Save(ctx context.Context, order *models.Order) error
}

package repositories

import (
	"context"
	"fmt"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
)

type orderRepository struct {
	orm contracts.ORM
}

var _ contracts.OrderRepository = (*orderRepository)(nil)

func NewOrderRepository(orm contracts.ORM) (contracts.OrderRepository, error) {
	if orm == nil {
		return nil, fmt.Errorf("new order repository: %w", ErrNilORM)
	}

	return &orderRepository{orm: orm}, nil
}

func (r *orderRepository) FindByAppmaxOrderID(_ context.Context, appmaxOrderID int) (*models.Order, error) {
	var order models.Order
	err := r.orm.Query().Where("appmax_order_id = ?", appmaxOrderID).First(&order)
	if err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}
	if order.ID == 0 {
		return nil, nil
	}
	return &order, nil
}

func (r *orderRepository) FindByAppmaxOrderIDAndInstallation(_ context.Context, appmaxOrderID int, installationID int64) (*models.Order, error) {
	var order models.Order
	err := r.orm.Query().
		Where("appmax_order_id = ? AND installation_id = ?", appmaxOrderID, installationID).
		First(&order)
	if err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}
	if order.ID == 0 {
		return nil, nil
	}
	return &order, nil
}

func (r *orderRepository) Create(_ context.Context, order *models.Order) error {
	return r.orm.Query().Create(order)
}

func (r *orderRepository) Save(_ context.Context, order *models.Order) error {
	return r.orm.Query().Save(order)
}

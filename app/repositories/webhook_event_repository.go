package repositories

import (
	"context"
	"fmt"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
)

type webhookEventRepository struct {
	orm contracts.ORM
}

var _ contracts.WebhookEventRepository = (*webhookEventRepository)(nil)

func NewWebhookEventRepository(orm contracts.ORM) (contracts.WebhookEventRepository, error) {
	if orm == nil {
		return nil, fmt.Errorf("new webhook event repository: %w", ErrNilORM)
	}

	return &webhookEventRepository{orm: orm}, nil
}

func (r *webhookEventRepository) Create(_ context.Context, event *models.WebhookEvent) error {
	return r.orm.Query().Create(event)
}

func (r *webhookEventRepository) Save(_ context.Context, event *models.WebhookEvent) error {
	return r.orm.Query().Save(event)
}

func (r *webhookEventRepository) FindProcessedDuplicate(_ context.Context, event string, appmaxOrderID int, excludeID int64) (*models.WebhookEvent, error) {
	var existing models.WebhookEvent
	err := r.orm.Query().
		Where("event = ? AND appmax_order_id = ? AND processed = ? AND id != ?",
			event, appmaxOrderID, true, excludeID).
		First(&existing)
	if err != nil {
		if isNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}
	if existing.ID == 0 {
		return nil, nil
	}
	return &existing, nil
}

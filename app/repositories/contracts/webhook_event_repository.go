package contracts

import (
	"context"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
)

type WebhookEventRepository interface {
	Create(ctx context.Context, event *models.WebhookEvent) error
	Save(ctx context.Context, event *models.WebhookEvent) error
	FindProcessedDuplicate(ctx context.Context, event string, appmaxOrderID int, excludeID int64) (*models.WebhookEvent, error)
}

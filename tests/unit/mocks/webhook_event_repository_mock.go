package mocks

import (
	"context"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
)

type MockWebhookEventRepository struct {
	CreateFunc                 func(ctx context.Context, event *models.WebhookEvent) error
	SaveFunc                   func(ctx context.Context, event *models.WebhookEvent) error
	FindProcessedDuplicateFunc func(ctx context.Context, event string, appmaxOrderID int, excludeID int64) (*models.WebhookEvent, error)
}

func (m *MockWebhookEventRepository) Create(ctx context.Context, event *models.WebhookEvent) error {
	return m.CreateFunc(ctx, event)
}

func (m *MockWebhookEventRepository) Save(ctx context.Context, event *models.WebhookEvent) error {
	return m.SaveFunc(ctx, event)
}

func (m *MockWebhookEventRepository) FindProcessedDuplicate(ctx context.Context, event string, appmaxOrderID int, excludeID int64) (*models.WebhookEvent, error) {
	return m.FindProcessedDuplicateFunc(ctx, event, appmaxOrderID, excludeID)
}

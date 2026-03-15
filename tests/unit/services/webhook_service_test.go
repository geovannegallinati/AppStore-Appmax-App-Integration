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

func intPtr(v int) *int { return &v }

func baseEventRepo(eventID int64) *mocks.MockWebhookEventRepository {
	return &mocks.MockWebhookEventRepository{
		CreateFunc: func(_ context.Context, event *models.WebhookEvent) error {
			event.ID = eventID
			return nil
		},
		SaveFunc: func(_ context.Context, _ *models.WebhookEvent) error {
			return nil
		},
		FindProcessedDuplicateFunc: func(_ context.Context, _ string, _ int, _ int64) (*models.WebhookEvent, error) {
			return nil, nil
		},
	}
}

func mustWebhookService(t *testing.T, eventRepo *mocks.MockWebhookEventRepository, orderRepo *mocks.MockOrderRepository) services.WebhookService {
	t.Helper()

	svc, err := services.NewWebhookService(eventRepo, orderRepo)
	require.NoError(t, err)

	return svc
}

func TestWebhookService_Handle_HappyPath(t *testing.T) {
	orderID := 99
	savedStatus := ""
	markedProcessed := false

	eventRepo := baseEventRepo(1)
	eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
		markedProcessed = event.Processed
		return nil
	}
	orderRepo := &mocks.MockOrderRepository{
		FindByAppmaxOrderIDFunc: func(_ context.Context, id int) (*models.Order, error) {
			return &models.Order{ID: 10, AppmaxOrderID: id, Status: "pendente"}, nil
		},
		SaveFunc: func(_ context.Context, order *models.Order) error {
			savedStatus = order.Status
			return nil
		},
	}

	svc := mustWebhookService(t, eventRepo, orderRepo)
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:     "order_paid",
		EventType: "order",
		OrderID:   &orderID,
		Payload:   models.JSONMap{"order_id": orderID},
	})

	require.NoError(t, err)
	assert.False(t, result.AlreadyProcessed)
	assert.Equal(t, "aprovado", savedStatus)
	assert.True(t, markedProcessed)
}

func TestWebhookService_Handle_DuplicateDetected(t *testing.T) {
	orderID := 99

	eventRepo := baseEventRepo(2)
	eventRepo.FindProcessedDuplicateFunc = func(_ context.Context, _ string, _ int, _ int64) (*models.WebhookEvent, error) {
		return &models.WebhookEvent{ID: 1, Processed: true}, nil
	}

	svc := mustWebhookService(t, eventRepo, &mocks.MockOrderRepository{})
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "order_paid",
		OrderID: &orderID,
		Payload: models.JSONMap{},
	})

	require.NoError(t, err)
	assert.True(t, result.AlreadyProcessed)
}

func TestWebhookService_Handle_UnknownEvent(t *testing.T) {
	orderID := 99
	markedProcessed := false

	eventRepo := baseEventRepo(3)
	eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
		markedProcessed = event.Processed
		return nil
	}

	svc := mustWebhookService(t, eventRepo, &mocks.MockOrderRepository{})
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "unknown_event_type",
		OrderID: &orderID,
		Payload: models.JSONMap{},
	})

	require.NoError(t, err)
	assert.False(t, result.AlreadyProcessed)
	assert.True(t, markedProcessed)
}

func TestWebhookService_Handle_NilOrderID(t *testing.T) {
	markedProcessed := false

	eventRepo := baseEventRepo(4)
	eventRepo.FindProcessedDuplicateFunc = func(_ context.Context, _ string, _ int, _ int64) (*models.WebhookEvent, error) {
		return nil, nil
	}
	eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
		markedProcessed = event.Processed
		return nil
	}

	svc := mustWebhookService(t, eventRepo, &mocks.MockOrderRepository{})
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "order_paid",
		OrderID: nil,
		Payload: models.JSONMap{},
	})

	require.NoError(t, err)
	assert.False(t, result.AlreadyProcessed)
	assert.True(t, markedProcessed)
}

func TestWebhookService_Handle_OrderNotFound(t *testing.T) {
	orderID := 99
	markedProcessed := false

	eventRepo := baseEventRepo(5)
	eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
		markedProcessed = event.Processed
		return nil
	}
	orderRepo := &mocks.MockOrderRepository{
		FindByAppmaxOrderIDFunc: func(_ context.Context, _ int) (*models.Order, error) {
			return nil, nil
		},
	}

	svc := mustWebhookService(t, eventRepo, orderRepo)
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "order_paid",
		OrderID: &orderID,
		Payload: models.JSONMap{},
	})

	require.NoError(t, err)
	assert.False(t, result.AlreadyProcessed)
	assert.True(t, markedProcessed)
}

func TestWebhookService_Handle_OrderSaveFails(t *testing.T) {
	orderID := 99
	saveErr := errors.New("db full")

	eventRepo := baseEventRepo(6)
	eventRepo.SaveFunc = func(_ context.Context, _ *models.WebhookEvent) error {
		return nil
	}
	orderRepo := &mocks.MockOrderRepository{
		FindByAppmaxOrderIDFunc: func(_ context.Context, id int) (*models.Order, error) {
			return &models.Order{ID: 10, AppmaxOrderID: id}, nil
		},
		SaveFunc: func(_ context.Context, _ *models.Order) error {
			return saveErr
		},
	}

	svc := mustWebhookService(t, eventRepo, orderRepo)
	_, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "order_paid",
		OrderID: &orderID,
		Payload: models.JSONMap{},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "db full")
}

func TestWebhookService_Handle_EventCreateFails(t *testing.T) {
	orderID := 99
	createErr := errors.New("insert failed")

	eventRepo := &mocks.MockWebhookEventRepository{
		CreateFunc: func(_ context.Context, _ *models.WebhookEvent) error {
			return createErr
		},
	}

	svc := mustWebhookService(t, eventRepo, &mocks.MockOrderRepository{})
	_, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "order_paid",
		OrderID: &orderID,
		Payload: models.JSONMap{},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "insert failed")
}

func TestWebhookServiceConstructor_RejectsNilDependency(t *testing.T) {
	svc, err := services.NewWebhookService(nil, &mocks.MockOrderRepository{})

	require.Error(t, err)
	assert.Nil(t, svc)
	assert.ErrorIs(t, err, services.ErrNilDependency)
}

func TestWebhookServiceConstructor_Success(t *testing.T) {
	svc, err := services.NewWebhookService(&mocks.MockWebhookEventRepository{}, &mocks.MockOrderRepository{})

	require.NoError(t, err)
	assert.NotNil(t, svc)
}

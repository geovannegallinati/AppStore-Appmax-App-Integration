package services

import (
	"context"
	"fmt"
	"time"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/models"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/repositories/contracts"
)

var webhookStatusMap = map[string]string{
	"order_authorized":              "autorizado",
	"order_authorized_with_delay":   "autorizado",
	"order_approved":                "aprovado",
	"order_billet_created":          "pendente",
	"order_paid":                    "aprovado",
	"order_pending_integration":     "pendente_integracao",
	"order_refund":                  "estornado",
	"order_pix_created":             "pendente",
	"order_paid_by_pix":             "aprovado",
	"order_pix_expired":             "cancelado",
	"order_integrated":              "integrado",
	"order_billet_overdue":          "cancelado",
	"order_chargeback_in_treatment": "chargeback_em_tratativa",
	"order_up_sold":                 "aprovado",
	"payment_not_authorized":        "cancelado",
	"payment_authorized_with_delay": "autorizado",
	"split_orders":                  "aprovado",
}

var knownNoOpEvents = map[string]bool{
	"customer_created":         true,
	"customer_interested":      true,
	"customer_contacted":       true,
	"subscription_cancelation": true,
	"subscription_delayed":     true,
}

type WebhookInput struct {
	Event     string
	EventType string
	OrderID   *int
	Payload   models.JSONMap
}

type WebhookResult struct {
	AlreadyProcessed bool
}

type WebhookService interface {
	Handle(ctx context.Context, input WebhookInput) (WebhookResult, error)
}

type webhookService struct {
	eventRepo contracts.WebhookEventRepository
	orderRepo contracts.OrderRepository
	logger    Logger
}

var _ WebhookService = (*webhookService)(nil)

func NewWebhookService(eventRepo contracts.WebhookEventRepository, orderRepo contracts.OrderRepository, logger ...Logger) (WebhookService, error) {
	if eventRepo == nil || orderRepo == nil {
		return nil, fmt.Errorf("new webhook service: %w", ErrNilDependency)
	}

	var selectedLogger Logger = noopLogger{}
	if len(logger) > 0 && logger[0] != nil {
		selectedLogger = logger[0]
	}

	return &webhookService{
		eventRepo: eventRepo,
		orderRepo: orderRepo,
		logger:    selectedLogger,
	}, nil
}

func (s *webhookService) Handle(ctx context.Context, input WebhookInput) (WebhookResult, error) {
	event := &models.WebhookEvent{
		Event:         input.Event,
		EventType:     input.EventType,
		AppmaxOrderID: input.OrderID,
		Payload:       input.Payload,
	}
	if err := s.eventRepo.Create(ctx, event); err != nil {
		return WebhookResult{}, fmt.Errorf("webhook handle: persist event: %w", err)
	}

	if input.OrderID != nil {
		dup, err := s.eventRepo.FindProcessedDuplicate(ctx, input.Event, *input.OrderID, event.ID)
		if err != nil {
			return WebhookResult{}, fmt.Errorf("webhook handle: dedup check: %w", err)
		}
		if dup != nil {
			return WebhookResult{AlreadyProcessed: true}, nil
		}
	}

	if knownNoOpEvents[input.Event] {
		return WebhookResult{}, s.markProcessed(ctx, event, "")
	}

	newStatus, hasMappedStatus := webhookStatusMap[input.Event]
	if !hasMappedStatus || input.OrderID == nil {
		s.logger.Warningf("webhook_service: event %s has no status mapping or no order_id, marking processed", input.Event)
		return WebhookResult{}, s.markProcessed(ctx, event, "")
	}

	order, err := s.orderRepo.FindByAppmaxOrderID(ctx, *input.OrderID)
	if err != nil {
		return WebhookResult{}, fmt.Errorf("webhook handle: find order: %w", err)
	}
	if order == nil {
		s.logger.Warningf("webhook_service: order %d not found for event %s", *input.OrderID, input.Event)
		return WebhookResult{}, s.markProcessed(ctx, event, "")
	}

	order.Status = newStatus
	if err := s.orderRepo.Save(ctx, order); err != nil {
		_ = s.markProcessed(ctx, event, err.Error())
		return WebhookResult{}, fmt.Errorf("webhook handle: update order status: %w", err)
	}

	return WebhookResult{}, s.markProcessed(ctx, event, "")
}

func (s *webhookService) markProcessed(ctx context.Context, event *models.WebhookEvent, errMsg string) error {
	now := time.Now()
	event.Processed = true
	event.ProcessedAt = &now
	event.ErrorMessage = errMsg
	if err := s.eventRepo.Save(ctx, event); err != nil {
		s.logger.Errorf("webhook_service: failed to mark event %d as processed: %v", event.ID, err)
		return err
	}
	return nil
}

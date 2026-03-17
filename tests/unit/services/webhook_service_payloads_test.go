package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/models"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/tests/unit/mocks"
)

func orderData(orderID int, appmaxStatus, paymentType string) map[string]any {
	return map[string]any{
		"id":                        orderID,
		"customer_id":               1,
		"total_products":            265,
		"status":                    appmaxStatus,
		"freight_value":             2.48,
		"freight_type":              "PAC",
		"payment_type":              paymentType,
		"card_brand":                nil,
		"card_details":              nil,
		"partner_total":             nil,
		"partner_affiliate_total":   nil,
		"pix_creation_date":         nil,
		"pix_expiration_date":       nil,
		"pix_emv":                   nil,
		"pix_ref":                   nil,
		"pix_qrcode":                nil,
		"pix_end_to_end_id":         nil,
		"billet_date_overdue":       "",
		"billet_url":                nil,
		"billet_digitable_line":     nil,
		"order_billet_payment_code": nil,
		"installments":              nil,
		"paid_at":                   nil,
		"refunded_at":               nil,
		"integrated_at":             nil,
		"created_at":                "2026-03-17 09:25:13",
		"discount":                  nil,
		"interest":                  nil,
		"upsell_order_id":           nil,
		"origin":                    "Site",
		"seller_name":               nil,
		"total":                     267.48,
		"total_refunded":            0,
		"refund_amount":             0,
		"full_payment_amount":       "267.48",
		"pix_payment_link":          "",
		"payment_link_id":           nil,
		"invoice_id":                nil,
		"issuer_message":            nil,
		"bundles":                   []any{},
		"products":                  []any{},
		"visit":                     []any{},
		"company_name":              "DemoAppGeovanne",
		"company_cnpj":              "66535306046",
		"company_email":             "geovanne.gallinati@teste.com",
		"co_production_commission":  nil,
		"affiliate_comission":       []any{},
		"traffic_description":       nil,
		"customer": map[string]any{
			"id":                        1,
			"site_id":                   1470,
			"firstname":                 "Pâmela",
			"lastname":                  "Furtado",
			"email":                     "daniel15@yahoo.com",
			"telephone":                 "3728675986",
			"postcode":                  "78448688",
			"address_street":            "Rua Verdugo",
			"address_street_number":     "3",
			"address_street_complement": "11º Andar",
			"address_street_district":   "laboriosam",
			"address_city":              "São Emanuel do Sul",
			"address_state":             "TO",
			"document_number":           "82132210799",
			"created_at":                "2026-03-17 09:25:13",
			"visited_url":               "http://example.com/test-product",
			"uf":                        "TO",
			"fullname":                  "Pâmela Furtado",
		},
	}
}

func buildOrderPayload(event, eventType string, orderID int, appmaxStatus, paymentType string) models.JSONMap {
	return models.JSONMap{
		"event":      event,
		"event_type": eventType,
		"data":       orderData(orderID, appmaxStatus, paymentType),
	}
}

func pixOrderData(orderID int, appmaxStatus, pixLink string) map[string]any {
	d := orderData(orderID, appmaxStatus, "Pix")
	d["pix_payment_link"] = pixLink
	return d
}

func buildCustomerPayload(event string, customerID int) models.JSONMap {
	return models.JSONMap{
		"event":      event,
		"event_type": "",
		"data": map[string]any{
			"id":                        customerID,
			"site_id":                   1470,
			"firstname":                 "Laura",
			"lastname":                  "Montenegro",
			"email":                     "maximo33@gmail.com",
			"telephone":                 "3526001145",
			"postcode":                  "97120572",
			"address_street":            "Avenida Santiago",
			"address_street_number":     "55",
			"address_street_complement": "Bloco C",
			"address_street_district":   "qui",
			"address_city":              "Santa Daniel",
			"address_state":             "SP",
			"document_number":           "57135108914",
			"created_at":                "2026-03-17 09:19:46",
			"visited_url":               "",
			"uf":                        "SP",
			"fullname":                  "Laura Montenegro",
			"interested_bundle":         []any{},
		},
	}
}

func buildSubscriptionPayload(event string, customerID int) models.JSONMap {
	return models.JSONMap{
		"event":      event,
		"event_type": "",
		"data": map[string]any{
			"id":                        customerID,
			"site_id":                   1470,
			"firstname":                 "Noelí",
			"lastname":                  "Guerra",
			"email":                     "felipe22@galhardo.biz",
			"telephone":                 "5191913915",
			"postcode":                  "91642854",
			"address_street":            "Travessa Mendes",
			"address_street_number":     "39332",
			"address_street_complement": "4º Andar",
			"address_street_district":   "voluptatum",
			"address_city":              "Noel do Norte",
			"address_state":             "SE",
			"document_number":           "75798407829",
			"created_at":                "2026-03-17 09:26:35",
			"visited_url":               "",
			"uf":                        "SE",
			"fullname":                  "Noelí Guerra",
			"interested_bundle":         []any{},
			"subscription": map[string]any{
				"id":           nil,
				"installments": nil,
				"total":        nil,
				"created_at":   "2026-03-17 09:26:35",
				"cancelled_at": nil,
				"bundle": map[string]any{
					"id":              1,
					"name":            "My product 1",
					"description":     "",
					"production_cost": "R$ 0,00",
					"identifier":      nil,
					"products": []any{
						map[string]any{
							"id":          1,
							"sku":         "9000010",
							"name":        "Livro de receitas",
							"description": "",
							"price":       "62.00",
							"quantity":    1,
							"image":       "https://gateway-boleto-homolog-appmax.s3-sa-east-1.amazonaws.com/",
							"external_id": nil,
						},
					},
				},
			},
		},
	}
}

// TestWebhookService_MappedEvents_UpdatesOrderStatus verifies that every event
// in webhookStatusMap correctly updates the order to its expected status.
func TestWebhookService_MappedEvents_UpdatesOrderStatus(t *testing.T) {
	tests := []struct {
		event       string
		wantStatus  string
		paymentType string
	}{
		{"order_authorized", "autorizado", "CreditCard"},
		{"order_authorized_with_delay", "autorizado", "CreditCard"},
		{"order_approved", "aprovado", "CreditCard"},
		{"order_billet_created", "pendente", "Boleto"},
		{"order_paid", "aprovado", "CreditCard"},
		{"order_pending_integration", "pendente_integracao", "CreditCard"},
		{"order_refund", "estornado", "CreditCard"},
		{"order_pix_created", "pendente", "Pix"},
		{"order_paid_by_pix", "aprovado", "Pix"},
		{"order_pix_expired", "cancelado", "Pix"},
		{"order_integrated", "integrado", "CreditCard"},
		{"order_billet_overdue", "cancelado", "Boleto"},
		{"order_chargeback_in_treatment", "chargeback_em_tratativa", "CreditCard"},
		{"order_up_sold", "aprovado", "CreditCard"},
		{"payment_not_authorized", "cancelado", "CreditCard"},
		{"payment_authorized_with_delay", "autorizado", "CreditCard"},
		{"split_orders", "aprovado", "CreditCard"},
	}

	for _, tt := range tests {
		t.Run(tt.event, func(t *testing.T) {
			orderID := 1
			savedStatus := ""
			markedProcessed := false

			eventRepo := baseEventRepo(10)
			eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
				markedProcessed = event.Processed
				return nil
			}
			orderRepo := &mocks.MockOrderRepository{
				FindByAppmaxOrderIDFunc: func(_ context.Context, id int) (*models.Order, error) {
					return &models.Order{ID: 5, AppmaxOrderID: id, Status: "pendente"}, nil
				},
				SaveFunc: func(_ context.Context, order *models.Order) error {
					savedStatus = order.Status
					return nil
				},
			}

			svc := mustWebhookService(t, eventRepo, orderRepo)
			result, err := svc.Handle(context.Background(), services.WebhookInput{
				Event:     tt.event,
				EventType: "order",
				OrderID:   &orderID,
				Payload:   buildOrderPayload(tt.event, "order", orderID, "pendente", tt.paymentType),
			})

			require.NoError(t, err)
			assert.False(t, result.AlreadyProcessed)
			assert.Equal(t, tt.wantStatus, savedStatus)
			assert.True(t, markedProcessed)
		})
	}
}

// TestWebhookService_KnownNoOpEvents_MarkedProcessedWithoutOrderLookup verifies that
// events in knownNoOpEvents are processed without touching the order repository.
func TestWebhookService_KnownNoOpEvents_MarkedProcessedWithoutOrderLookup(t *testing.T) {
	tests := []struct {
		event   string
		payload models.JSONMap
	}{
		{"customer_created", buildCustomerPayload("customer_created", 1)},
		{"customer_interested", buildCustomerPayload("customer_interested", 1)},
		{"customer_contacted", buildCustomerPayload("customer_contacted", 1)},
		{"subscription_cancelation", buildSubscriptionPayload("subscription_cancelation", 1)},
		{"subscription_delayed", buildSubscriptionPayload("subscription_delayed", 1)},
	}

	for _, tt := range tests {
		t.Run(tt.event, func(t *testing.T) {
			markedProcessed := false
			orderLookupCalled := false

			eventRepo := baseEventRepo(20)
			eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
				markedProcessed = event.Processed
				return nil
			}
			orderRepo := &mocks.MockOrderRepository{
				FindByAppmaxOrderIDFunc: func(_ context.Context, _ int) (*models.Order, error) {
					orderLookupCalled = true
					return nil, nil
				},
			}

			svc := mustWebhookService(t, eventRepo, orderRepo)
			result, err := svc.Handle(context.Background(), services.WebhookInput{
				Event:   tt.event,
				OrderID: nil,
				Payload: tt.payload,
			})

			require.NoError(t, err)
			assert.False(t, result.AlreadyProcessed)
			assert.True(t, markedProcessed)
			assert.False(t, orderLookupCalled)
		})
	}
}

// TestWebhookService_PascalCaseMappedEvents_UpdatesOrderStatus verifies that every
// PascalCase event in webhookStatusMap correctly updates the order to its expected status.
func TestWebhookService_PascalCaseMappedEvents_UpdatesOrderStatus(t *testing.T) {
	tests := []struct {
		event       string
		wantStatus  string
		paymentType string
	}{
		{"OrderAuthorized", "autorizado", "CreditCard"},
		{"OrderApproved", "aprovado", "CreditCard"},
		{"OrderBilletCreated", "pendente", "Boleto"},
		{"OrderPaid", "aprovado", "CreditCard"},
		{"OrderPendingIntegration", "pendente_integracao", "CreditCard"},
		{"OrderRefund", "estornado", "CreditCard"},
		{"OrderPixCreated", "pendente", "Pix"},
		{"OrderPaidByPix", "aprovado", "Pix"},
		{"OrderPixExpired", "cancelado", "Pix"},
		{"OrderIntegrated", "integrado", "CreditCard"},
		{"OrderBilletOverdue", "cancelado", "Boleto"},
		{"OrderChargeBackInTreatment", "chargeback_em_tratativa", "CreditCard"},
		{"OrderUpSold", "aprovado", "CreditCard"},
		{"CreatedSubscription", "aprovado", "CreditCard"},
	}

	for _, tt := range tests {
		t.Run(tt.event, func(t *testing.T) {
			orderID := 1
			savedStatus := ""
			markedProcessed := false

			eventRepo := baseEventRepo(30)
			eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
				markedProcessed = event.Processed
				return nil
			}
			orderRepo := &mocks.MockOrderRepository{
				FindByAppmaxOrderIDFunc: func(_ context.Context, id int) (*models.Order, error) {
					return &models.Order{ID: 5, AppmaxOrderID: id, Status: "pendente"}, nil
				},
				SaveFunc: func(_ context.Context, order *models.Order) error {
					savedStatus = order.Status
					return nil
				},
			}

			svc := mustWebhookService(t, eventRepo, orderRepo)
			result, err := svc.Handle(context.Background(), services.WebhookInput{
				Event:   tt.event,
				OrderID: &orderID,
				Payload: buildOrderPayload(tt.event, "", orderID, "pendente", tt.paymentType),
			})

			require.NoError(t, err)
			assert.False(t, result.AlreadyProcessed)
			assert.Equal(t, tt.wantStatus, savedStatus)
			assert.True(t, markedProcessed)
		})
	}
}

// TestWebhookService_PascalCaseUnmappedEvents_MarkedProcessedWithoutStatusUpdate verifies
// that PascalCase events not present in webhookStatusMap are persisted and marked
// processed without updating order status.
func TestWebhookService_PascalCaseUnmappedEvents_MarkedProcessedWithoutStatusUpdate(t *testing.T) {
	tests := []struct {
		event       string
		paymentType string
	}{
		{"OrderPartialRefund", "CreditCard"},
		{"OrderChargeBackGain", "CreditCard"},
		{"ChargeFailed", "CreditCard"},
		{"ChargeSuccess", "CreditCard"},
		{"CanceledSubscription", "CreditCard"},
	}

	for _, tt := range tests {
		t.Run(tt.event, func(t *testing.T) {
			orderID := 1
			orderStatusUpdated := false
			markedProcessed := false

			eventRepo := baseEventRepo(31)
			eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
				markedProcessed = event.Processed
				return nil
			}
			orderRepo := &mocks.MockOrderRepository{
				FindByAppmaxOrderIDFunc: func(_ context.Context, id int) (*models.Order, error) {
					return &models.Order{ID: 5, AppmaxOrderID: id, Status: "pendente"}, nil
				},
				SaveFunc: func(_ context.Context, _ *models.Order) error {
					orderStatusUpdated = true
					return nil
				},
			}

			svc := mustWebhookService(t, eventRepo, orderRepo)
			result, err := svc.Handle(context.Background(), services.WebhookInput{
				Event:   tt.event,
				OrderID: &orderID,
				Payload: buildOrderPayload(tt.event, "", orderID, "pendente", tt.paymentType),
			})

			require.NoError(t, err)
			assert.False(t, result.AlreadyProcessed)
			assert.False(t, orderStatusUpdated)
			assert.True(t, markedProcessed)
		})
	}
}

// TestWebhookService_PascalCaseCustomerNoOpEvents_MarkedProcessedWithoutOrderLookup verifies
// that PascalCase customer and subscription events are handled as no-ops without touching
// the order repository.
func TestWebhookService_PascalCaseCustomerNoOpEvents_MarkedProcessedWithoutOrderLookup(t *testing.T) {
	tests := []struct {
		event   string
		payload models.JSONMap
	}{
		{"CustomerCreated", buildCustomerPayload("CustomerCreated", 1)},
		{"CustomerInterested", buildCustomerPayload("CustomerInterested", 1)},
		{"CustomerContacted", buildCustomerPayload("CustomerContacted", 1)},
		{"SubscriptionCancellationEvent", buildSubscriptionPayload("SubscriptionCancellationEvent", 1)},
		{"SubscriptionDelayedEvent", buildSubscriptionPayload("SubscriptionDelayedEvent", 1)},
	}

	for _, tt := range tests {
		t.Run(tt.event, func(t *testing.T) {
			markedProcessed := false
			orderLookupCalled := false

			eventRepo := baseEventRepo(40)
			eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
				markedProcessed = event.Processed
				return nil
			}
			orderRepo := &mocks.MockOrderRepository{
				FindByAppmaxOrderIDFunc: func(_ context.Context, _ int) (*models.Order, error) {
					orderLookupCalled = true
					return nil, nil
				},
			}

			svc := mustWebhookService(t, eventRepo, orderRepo)
			result, err := svc.Handle(context.Background(), services.WebhookInput{
				Event:   tt.event,
				OrderID: nil,
				Payload: tt.payload,
			})

			require.NoError(t, err)
			assert.False(t, result.AlreadyProcessed)
			assert.True(t, markedProcessed)
			assert.False(t, orderLookupCalled)
		})
	}
}

// TestWebhookService_SubscriptionCancellationEvent_WithSubscriptionObject verifies
// that the nested subscription object in the payload is persisted without errors.
func TestWebhookService_SubscriptionCancellationEvent_WithSubscriptionObject(t *testing.T) {
	markedProcessed := false

	eventRepo := baseEventRepo(50)
	eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
		markedProcessed = event.Processed
		return nil
	}

	svc := mustWebhookService(t, eventRepo, &mocks.MockOrderRepository{})
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "SubscriptionCancellationEvent",
		OrderID: nil,
		Payload: buildSubscriptionPayload("SubscriptionCancellationEvent", 1),
	})

	require.NoError(t, err)
	assert.False(t, result.AlreadyProcessed)
	assert.True(t, markedProcessed)
}

// TestWebhookService_SubscriptionDelayedEvent_WithSubscriptionObject verifies
// that the nested subscription object in the payload is persisted without errors.
func TestWebhookService_SubscriptionDelayedEvent_WithSubscriptionObject(t *testing.T) {
	markedProcessed := false

	eventRepo := baseEventRepo(51)
	eventRepo.SaveFunc = func(_ context.Context, event *models.WebhookEvent) error {
		markedProcessed = event.Processed
		return nil
	}

	svc := mustWebhookService(t, eventRepo, &mocks.MockOrderRepository{})
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "SubscriptionDelayedEvent",
		OrderID: nil,
		Payload: buildSubscriptionPayload("SubscriptionDelayedEvent", 1),
	})

	require.NoError(t, err)
	assert.False(t, result.AlreadyProcessed)
	assert.True(t, markedProcessed)
}

// TestWebhookService_OrderPixCreated_PixPaymentLinkInPayload verifies that
// the snake_case event (order_pix_created) sets the order status to "pendente"
// and the pix_payment_link field from the real Appmax payload is stored correctly.
func TestWebhookService_OrderPixCreated_PixPaymentLinkInPayload(t *testing.T) {
	orderID := 1
	savedStatus := ""

	eventRepo := baseEventRepo(60)
	eventRepo.SaveFunc = func(_ context.Context, _ *models.WebhookEvent) error { return nil }

	orderRepo := &mocks.MockOrderRepository{
		FindByAppmaxOrderIDFunc: func(_ context.Context, id int) (*models.Order, error) {
			return &models.Order{ID: 5, AppmaxOrderID: id, Status: "novo"}, nil
		},
		SaveFunc: func(_ context.Context, order *models.Order) error {
			savedStatus = order.Status
			return nil
		},
	}

	pixLink := "https://breakingcode.sandboxappmax.com.br/show-pix/1"
	payload := models.JSONMap{
		"event":      "order_pix_created",
		"event_type": "",
		"data":       pixOrderData(orderID, "pendente", pixLink),
	}

	svc := mustWebhookService(t, eventRepo, orderRepo)
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "order_pix_created",
		OrderID: &orderID,
		Payload: payload,
	})

	require.NoError(t, err)
	assert.False(t, result.AlreadyProcessed)
	assert.Equal(t, "pendente", savedStatus)
	assert.Equal(t, pixLink, payload["data"].(map[string]any)["pix_payment_link"])
}

// TestWebhookService_OrderPaidByPix_PixPaymentLinkInPayload verifies that
// order_paid_by_pix sets status to "aprovado" and the pix_payment_link is present.
func TestWebhookService_OrderPaidByPix_PixPaymentLinkInPayload(t *testing.T) {
	orderID := 1
	savedStatus := ""

	eventRepo := baseEventRepo(61)
	eventRepo.SaveFunc = func(_ context.Context, _ *models.WebhookEvent) error { return nil }

	orderRepo := &mocks.MockOrderRepository{
		FindByAppmaxOrderIDFunc: func(_ context.Context, id int) (*models.Order, error) {
			return &models.Order{ID: 5, AppmaxOrderID: id, Status: "pendente"}, nil
		},
		SaveFunc: func(_ context.Context, order *models.Order) error {
			savedStatus = order.Status
			return nil
		},
	}

	pixLink := "https://breakingcode.sandboxappmax.com.br/show-pix/1"
	payload := models.JSONMap{
		"event":      "order_paid_by_pix",
		"event_type": "",
		"data":       pixOrderData(orderID, "aprovado", pixLink),
	}

	svc := mustWebhookService(t, eventRepo, orderRepo)
	result, err := svc.Handle(context.Background(), services.WebhookInput{
		Event:   "order_paid_by_pix",
		OrderID: &orderID,
		Payload: payload,
	})

	require.NoError(t, err)
	assert.False(t, result.AlreadyProcessed)
	assert.Equal(t, "aprovado", savedStatus)
	assert.Equal(t, pixLink, payload["data"].(map[string]any)["pix_payment_link"])
}

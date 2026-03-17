package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/requests"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/responses"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/models"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
)

type WebhookController struct {
	webhookSvc services.WebhookService
}

func NewWebhookController(webhookSvc services.WebhookService) (*WebhookController, error) {
	if webhookSvc == nil {
		return nil, fmt.Errorf("new webhook controller: %w", ErrNilDependency)
	}

	return &WebhookController{webhookSvc: webhookSvc}, nil
}

func (c *WebhookController) Handle(ctx http.Context) http.Response {
	var envelope requests.WebhookEnvelopeRequest
	if err := ctx.Request().Bind(&envelope); err != nil {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	var data requests.WebhookDataRequest
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		// keep backward compatibility with payload variations while extracting order_id
	}
	orderID := data.ExtractOrderID()
	model := data.DetectModel(envelope.EventType)
	logIncomingWebhook(envelope, orderID, model)

	payload := models.JSONMap{
		"event":      envelope.Event,
		"event_type": envelope.EventType,
		"data":       json.RawMessage(envelope.Data),
	}

	result, err := c.webhookSvc.Handle(ctx.Context(), services.WebhookInput{
		Event:     envelope.Event,
		EventType: envelope.EventType,
		OrderID:   orderID,
		Payload:   payload,
	})
	if err != nil {
		facades.Log().Errorf("webhook_controller: handle failed for event %s: %v", envelope.Event, err)
		return ctx.Response().Json(500, responses.MessageResponse{Message: "internal server error"})
	}

	if result.AlreadyProcessed {
		return ctx.Response().Json(200, responses.MessageResponse{Message: "already processed"})
	}
	return ctx.Response().Json(200, responses.MessageResponse{Message: "ok"})
}

func logIncomingWebhook(envelope requests.WebhookEnvelopeRequest, orderID *int, model string) {
	orderIDValue := "null"
	if orderID != nil {
		orderIDValue = fmt.Sprintf("%d", *orderID)
	}

	raw, err := json.Marshal(envelope)
	if err != nil {
		facades.Log().Debugf(
			"webhook_controller: received event=%s event_type=%s order_id=%s model=%s payload_unmarshalable=true payload_data=%s",
			envelope.Event,
			envelope.EventType,
			orderIDValue,
			model,
			strings.TrimSpace(string(envelope.Data)),
		)
		return
	}

	facades.Log().Debugf(
		"webhook_controller: received event=%s event_type=%s order_id=%s model=%s payload=%s",
		envelope.Event,
		envelope.EventType,
		orderIDValue,
		model,
		string(raw),
	)
}

package controllers_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/controllers"
	"github.com/stretchr/testify/assert"
)

type fakeUpstreamErr struct {
	status  int
	message string
}

func (e fakeUpstreamErr) Error() string {
	return fmt.Sprintf("unexpected status %d: {\"errors\":{\"message\":\"%s\"}}", e.status, e.message)
}

func (e fakeUpstreamErr) HTTPStatus() int {
	return e.status
}

func (e fakeUpstreamErr) UpstreamMessage() string {
	return e.message
}

func TestRefundErrorMessage_UsesUpstreamErrorsMessage(t *testing.T) {
	err := errors.New(`appmax refund: refund: unexpected status 404 CF-Ray=abc: {"errors":{"message":"Producer has no amount to realize this action"}}`)

	assert.Equal(t, "Producer has no amount to realize this action", controllers.RefundErrorMessage(err))
}

func TestRefundErrorMessage_UsesUpstreamMessage(t *testing.T) {
	err := errors.New(`appmax refund: refund: unexpected status 400: {"message":"Invalid request payload"}`)

	assert.Equal(t, "Invalid request payload", controllers.RefundErrorMessage(err))
}

func TestRefundErrorMessage_UsesFallbackWhenMessageUnavailable(t *testing.T) {
	err := errors.New("appmax refund: timeout")

	assert.Equal(t, controllers.RefundFailedMessage, controllers.RefundErrorMessage(err))
}

func TestUpstreamErrorMessage_UsesFallbackWhenNilError(t *testing.T) {
	assert.Equal(t, "fallback message", controllers.UpstreamErrorMessage(nil, "fallback message"))
}

func TestUpstreamErrorMessage_ExtractsNestedErrorsMessage(t *testing.T) {
	err := errors.New(`checkout upsell: appmax upsell: upsell: unexpected status 404 CF-Ray=xyz: {"errors":{"message":"Order not found."}}`)

	assert.Equal(t, "Order not found.", controllers.UpstreamErrorMessage(err, "upsell failed"))
}

func TestUpstreamErrorStatus_UsesTypedStatusFromWrappedError(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", fakeUpstreamErr{status: 404, message: "Order not found."})

	assert.Equal(t, 404, controllers.UpstreamErrorStatus(err, 502))
}

func TestUpstreamErrorStatus_ParsesStatusFromRawErrorMessage(t *testing.T) {
	err := errors.New(`appmax refund: refund: unexpected status 401 CF-Ray=abc: {"errors":{"message":"Unauthorized"}}`)

	assert.Equal(t, 401, controllers.UpstreamErrorStatus(err, 502))
}

func TestUpstreamErrorStatus_UsesFallbackWhenStatusUnavailable(t *testing.T) {
	err := errors.New("network timeout")

	assert.Equal(t, 502, controllers.UpstreamErrorStatus(err, 502))
}

func TestUpstreamErrorMessage_UsesTypedMessageFromWrappedError(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", fakeUpstreamErr{status: 403, message: "Forbidden"})

	assert.Equal(t, "Forbidden", controllers.UpstreamErrorMessage(err, "fallback"))
}

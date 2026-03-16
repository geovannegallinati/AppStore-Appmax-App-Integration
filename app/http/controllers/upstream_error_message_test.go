package controllers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefundErrorMessage_UsesUpstreamErrorsMessage(t *testing.T) {
	err := errors.New(`appmax refund: refund: unexpected status 404 CF-Ray=abc: {"errors":{"message":"Producer has no amount to realize this action"}}`)

	assert.Equal(t, "Producer has no amount to realize this action", refundErrorMessage(err))
}

func TestRefundErrorMessage_UsesUpstreamMessage(t *testing.T) {
	err := errors.New(`appmax refund: refund: unexpected status 400: {"message":"Invalid request payload"}`)

	assert.Equal(t, "Invalid request payload", refundErrorMessage(err))
}

func TestRefundErrorMessage_UsesFallbackWhenMessageUnavailable(t *testing.T) {
	err := errors.New("appmax refund: timeout")

	assert.Equal(t, refundFailedMessage, refundErrorMessage(err))
}

func TestUpstreamErrorMessage_UsesFallbackWhenNilError(t *testing.T) {
	assert.Equal(t, "fallback message", upstreamErrorMessage(nil, "fallback message"))
}

func TestUpstreamErrorMessage_ExtractsNestedErrorsMessage(t *testing.T) {
	err := errors.New(`checkout upsell: appmax upsell: upsell: unexpected status 404 CF-Ray=xyz: {"errors":{"message":"Order not found."}}`)

	assert.Equal(t, "Order not found.", upstreamErrorMessage(err, "upsell failed"))
}

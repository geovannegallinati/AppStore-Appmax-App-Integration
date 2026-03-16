package controllers

import (
	"fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/middleware"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/responses"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/models"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
)

type MerchantAuthController struct {
	tokenManager services.TokenManager
}

func NewMerchantAuthController(tokenManager services.TokenManager) (*MerchantAuthController, error) {
	if tokenManager == nil {
		return nil, fmt.Errorf("new merchant auth controller: %w", ErrNilDependency)
	}

	return &MerchantAuthController{tokenManager: tokenManager}, nil
}

func merchantInstallationFromCtx(ctx http.Context) (*models.Installation, bool) {
	inst, ok := ctx.Value(middleware.InstallationContextKey).(*models.Installation)
	return inst, ok && inst != nil
}

func (c *MerchantAuthController) SyncToken(ctx http.Context) http.Response {
	inst, ok := merchantInstallationFromCtx(ctx)
	if !ok {
		return ctx.Response().Json(500, responses.MessageResponse{Message: "installation context missing"})
	}

	token, err := c.tokenManager.MerchantToken(ctx.Context(), inst)
	if err != nil {
		facades.Log().Errorf("merchant_auth_controller: merchant token fetch failed for key %s: %v", inst.ExternalKey, err)
		return ctx.Response().Json(502, responses.MessageResponse{Message: upstreamErrorMessage(err, "failed to fetch merchant token")})
	}

	return ctx.Response().Json(200, responses.MerchantTokenSyncResponse{
		ExternalKey:         inst.ExternalKey,
		ExternalID:          inst.ExternalID,
		MerchantClientID:    inst.MerchantClientID,
		MerchantBearerToken: token,
		TokenType:           "Bearer",
	})
}

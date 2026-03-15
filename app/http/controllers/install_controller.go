package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/geovanne-gallinati/AppStoreAppDemo/app/http/requests"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/http/responses"
	"github.com/geovanne-gallinati/AppStoreAppDemo/app/services"
)

const appmaxCallTimeout = 3 * time.Minute

const installStateTTL = time.Hour

type installState struct {
	AppID       string
	ExternalKey string
}

type InstallController struct {
	appmaxSvc    services.AppmaxService
	installSvc   services.InstallService
	adminURL     string
	appURL       string
	appIDUUID    string
	appIDNumeric string
}

func NewInstallController(appmaxSvc services.AppmaxService, installSvc services.InstallService, adminURL, appURL, appIDUUID, appIDNumeric string) (*InstallController, error) {
	if appmaxSvc == nil || installSvc == nil {
		return nil, fmt.Errorf("new install controller: %w", ErrNilDependency)
	}
	if strings.TrimSpace(adminURL) == "" || strings.TrimSpace(appURL) == "" || strings.TrimSpace(appIDUUID) == "" || strings.TrimSpace(appIDNumeric) == "" {
		return nil, fmt.Errorf("new install controller: %w", ErrInvalidConfig)
	}

	return &InstallController{appmaxSvc: appmaxSvc, installSvc: installSvc, adminURL: adminURL, appURL: appURL, appIDUUID: appIDUUID, appIDNumeric: appIDNumeric}, nil
}

func (c *InstallController) Start(ctx http.Context) http.Response {
	appID := ctx.Request().Query("app_id", "")
	if appID == "" {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "app_id is required"})
	}

	externalKey := ctx.Request().Query("external_key", "")
	if externalKey == "" {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "external_key is required"})
	}
	attemptExternalKey := buildAttemptExternalKey(externalKey)

	callbackURL := fmt.Sprintf("%s/integrations/appmax/callback/install", c.appURL)

	appmaxCtx, cancel := context.WithTimeout(context.Background(), appmaxCallTimeout)
	defer cancel()

	hash, err := c.appmaxSvc.Authorize(appmaxCtx, appID, attemptExternalKey, callbackURL)
	if err != nil {
		facades.Log().Errorf("install_controller: authorize failed for key %s (attempt %s): %v", externalKey, attemptExternalKey, err)
		return ctx.Response().Json(502, responses.MessageResponse{Message: "failed to initiate installation"})
	}

	stateJSON, err := json.Marshal(installState{AppID: appID, ExternalKey: attemptExternalKey})
	if err != nil {
		facades.Log().Errorf("install_controller: marshal state failed for hash %s: %v", hash, err)
		return ctx.Response().Json(500, responses.MessageResponse{Message: "internal server error"})
	}
	if err := facades.Cache().Put("install:"+hash, string(stateJSON), installStateTTL); err != nil {
		facades.Log().Errorf("install_controller: cache put failed for hash %s: %v", hash, err)
		return ctx.Response().Json(500, responses.MessageResponse{Message: "internal server error"})
	}

	redirectURL := fmt.Sprintf("%s/appstore/integration/%s", c.adminURL, hash)
	return ctx.Response().Redirect(302, redirectURL)
}

func (c *InstallController) CallbackGuide(ctx http.Context) http.Response {
	token := strings.TrimSpace(ctx.Request().Query("token", ""))
	if token == "" {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "token is required"})
	}

	stateJSON := facades.Cache().GetString("install:" + token)
	if stateJSON == "" {
		facades.Log().Debugf("install_controller: no cached state for token %s — installation will be confirmed via health check POST", token)
		return ctx.Response().Json(200, responses.MessageResponse{Message: "installation confirmed"})
	}
	facades.Cache().Forget("install:" + token)

	var state installState
	if err := json.Unmarshal([]byte(stateJSON), &state); err != nil {
		facades.Log().Errorf("install_controller: unmarshal state failed for token %s: %v", token, err)
		return ctx.Response().Json(500, responses.MessageResponse{Message: "internal server error"})
	}

	if state.AppID != c.appIDUUID {
		facades.Log().Errorf("install_controller: app_id mismatch in state for token %s: got %s", token, state.AppID)
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid app_id"})
	}

	appmaxCtx, cancel := context.WithTimeout(context.Background(), appmaxCallTimeout)
	defer cancel()

	clientID, clientSecret, err := c.appmaxSvc.GenerateMerchantCreds(appmaxCtx, token)
	if err != nil {
		facades.Log().Errorf("install_controller: generate merchant creds failed for token %s: %v", token, err)
		return ctx.Response().Json(502, responses.MessageResponse{Message: "failed to generate merchant credentials"})
	}

	inst, created, err := c.installSvc.Upsert(ctx.Context(), services.UpsertInstallationInput{
		AppID:                state.AppID,
		ExternalKey:          state.ExternalKey,
		MerchantClientID:     clientID,
		MerchantClientSecret: clientSecret,
	})
	if err != nil {
		facades.Log().Errorf("install_controller: upsert failed for key %s: %v", state.ExternalKey, err)
		return ctx.Response().Json(500, responses.MessageResponse{Message: "internal server error"})
	}

	_ = created
	return ctx.Response().Json(200, responses.InstallCallbackResponse{ExternalID: inst.ExternalID})
}

func (c *InstallController) Callback(ctx http.Context) http.Response {
	var body requests.InstallCallbackRequest
	if err := ctx.Request().Bind(&body); err != nil {
		facades.Log().Debugf("[install] healthcheck POST from %s — bind error: %v", ctx.Request().Ip(), err)
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid request body"})
	}

	facades.Log().Debugf("[install] healthcheck POST from %s — app_id=%s external_key=%s client_id=%s",
		ctx.Request().Ip(), body.AppID, body.ExternalKey, body.MerchantClientID)

	if body.AppID == "" || body.ExternalKey == "" || body.MerchantClientID == "" || body.MerchantClientSecret == "" {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "app_id, external_key, merchant_client_id and merchant_client_secret are required"})
	}

	if body.AppID != c.appIDNumeric {
		facades.Log().Errorf("install_controller: app_id mismatch for key %s: got %s", body.ExternalKey, body.AppID)
		return ctx.Response().Json(400, responses.MessageResponse{Message: "invalid app_id"})
	}

	inst, created, err := c.installSvc.Upsert(ctx.Context(), services.UpsertInstallationInput{
		AppID:                body.AppID,
		ExternalKey:          body.ExternalKey,
		MerchantClientID:     body.MerchantClientID,
		MerchantClientSecret: body.MerchantClientSecret,
	})
	if err != nil {
		facades.Log().Errorf("install_controller: upsert failed for key %s: %v", body.ExternalKey, err)
		return ctx.Response().Json(500, responses.MessageResponse{Message: "internal server error"})
	}

	_ = created
	return ctx.Response().Json(200, responses.InstallCallbackResponse{ExternalID: inst.ExternalID})
}

func buildAttemptExternalKey(origin string) string {
	base := strings.TrimSpace(origin)
	if base == "" {
		base = "install"
	}

	suffixBytes := make([]byte, 8)
	if _, err := rand.Read(suffixBytes); err != nil {
		return fmt.Sprintf("%s-%d", truncateExternalKeyBase(base, 235), time.Now().UnixNano())
	}

	suffix := hex.EncodeToString(suffixBytes)
	maxBaseLen := 255 - 1 - len(suffix)
	if maxBaseLen < 1 {
		maxBaseLen = 1
	}

	return fmt.Sprintf("%s-%s", truncateExternalKeyBase(base, maxBaseLen), suffix)
}

func truncateExternalKeyBase(base string, maxLen int) string {
	if maxLen < 1 {
		return "k"
	}
	if len(base) <= maxLen {
		return base
	}
	return base[:maxLen]
}

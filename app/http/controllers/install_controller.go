package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/requests"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/responses"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
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
	appID := strings.TrimSpace(ctx.Request().Query("app_id", ""))
	externalKey := strings.TrimSpace(ctx.Request().Query("external_key", ""))
	if appID == "" || externalKey == "" {
		if requestWantsHTML(ctx) {
			return c.InstallStartFrontend(ctx, appID, externalKey)
		}
	}

	if appID == "" {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "app_id is required"})
	}

	if externalKey == "" {
		return ctx.Response().Json(400, responses.MessageResponse{Message: "external_key is required"})
	}

	callbackURL := fmt.Sprintf("%s/integrations/appmax/callback/install", c.appURL)

	appmaxCtx, cancel := context.WithTimeout(context.Background(), appmaxCallTimeout)
	defer cancel()

	hash, err := c.appmaxSvc.Authorize(appmaxCtx, appID, externalKey, callbackURL)
	if err != nil {
		facades.Log().Errorf("install_controller: authorize failed for key %s: %v", externalKey, err)
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "failed to initiate installation")})
	}

	stateJSON, err := json.Marshal(installState{AppID: appID, ExternalKey: externalKey})
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
	wantsHTML := requestWantsHTML(ctx)

	token := strings.TrimSpace(ctx.Request().Query("token", ""))
	if token == "" {
		if wantsHTML {
			return NewHealthController(c.appURL).CallbackFrontend(ctx)
		}
		return ctx.Response().Json(400, responses.MessageResponse{Message: "token is required"})
	}

	stateJSON := facades.Cache().GetString("install:" + token)
	if stateJSON == "" {
		facades.Log().Debugf("install_controller: no cached state for token %s — installation will be confirmed via health check POST", token)
		if wantsHTML {
			return c.InstallCompletedFrontend(ctx)
		}
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
		return ctx.Response().Json(UpstreamErrorStatus(err, 502), responses.MessageResponse{Message: UpstreamErrorMessage(err, "failed to generate merchant credentials")})
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
	if wantsHTML {
		return c.InstallCompletedFrontend(ctx)
	}
	return ctx.Response().Json(200, responses.InstallCallbackResponse{ExternalID: inst.ExternalID})
}

func (c *InstallController) InstallCompletedFrontend(ctx http.Context) http.Response {
	return ctx.Response().View().Make("frontend/install_completed.tmpl", map[string]any{
		"Title": "AppMax Checkout Demo - Installation Completed",
	})
}

func (c *InstallController) InstallStartFrontend(ctx http.Context, appID, externalKey string) http.Response {
	if appID == "" {
		appID = c.appIDUUID
	}
	if externalKey == "" {
		externalKey = fmt.Sprintf("install-%d", time.Now().Unix())
	}

	baseURL := frontendBaseURL(ctx, c.appURL)
	page := frontendPageData{
		Title:               "AppMax Checkout Demo - Install",
		Badge:               "Install",
		Headline:            "Start integration installation",
		Message:             "Fill in the fields below and use one click to start /install/start.",
		Submessage:          "The install button opens a new tab and redirects automatically to Breaking Code.",
		ActiveRoute:         routeInstallStart,
		PageKind:            "install",
		Endpoints:           frontendEndpoints(baseURL, routeInstallStart),
		ShowInstallForm:     true,
		InstallFormAction:   routeInstallStart,
		DefaultInstallAppID: appID,
		DefaultExternalKey:  externalKey,
		Tips: []string{
			"The install button sends a GET to /install/start with app_id and external_key.",
			"If you call /install/start without Accept text/html, the API keeps its original JSON behavior.",
		},
	}

	return ctx.Response().View().Make("frontend/page.tmpl", page)
}

func (c *InstallController) Callback(ctx http.Context) http.Response {
	var rawBody []byte
	if origin := ctx.Request().Origin(); origin != nil && origin.Body != nil {
		rawBody, _ = io.ReadAll(origin.Body)
		origin.Body = io.NopCloser(bytes.NewReader(rawBody))
	}
	facades.Log().Infof("[install:debug] raw callback from %s — headers: %v | body: %s",
		ctx.Request().Ip(),
		ctx.Request().Headers(),
		string(rawBody))

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

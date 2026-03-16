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

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/requests"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/responses"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
)

const appmaxCallTimeout = 3 * time.Minute

const installStateTTL = time.Hour
const installCompletedFrontendTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>AppMax Checkout Demo - Installation Completed</title>
  <style>
    :root {
      --bg: #0d1b2a;
      --panel: #1b263b;
      --ok: #2ec4b6;
      --txt: #e0e1dd;
      --muted: #9aa6b2;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      display: grid;
      place-items: center;
      font-family: "Trebuchet MS", "Segoe UI", sans-serif;
      color: var(--txt);
      background:
        radial-gradient(circle at 20% 20%, #22375a 0%, transparent 45%),
        radial-gradient(circle at 80% 80%, #1e3a34 0%, transparent 40%),
        var(--bg);
      padding: 16px;
    }
    .card {
      width: min(92vw, 640px);
      border: 1px solid #34445e;
      background: linear-gradient(160deg, rgba(27,38,59,0.95), rgba(17,28,43,0.95));
      border-radius: 16px;
      padding: 28px 24px;
      box-shadow: 0 14px 40px rgba(0,0,0,0.35);
    }
    .badge {
      display: inline-block;
      background: rgba(46,196,182,0.14);
      border: 1px solid rgba(46,196,182,0.35);
      color: var(--ok);
      border-radius: 999px;
      padding: 6px 12px;
      font-size: 12px;
      letter-spacing: 0.04em;
      text-transform: uppercase;
      margin-bottom: 12px;
    }
    h1 {
      margin: 0 0 10px;
      font-size: clamp(24px, 3.2vw, 34px);
      line-height: 1.2;
    }
    p {
      margin: 0 0 12px;
      color: var(--muted);
      line-height: 1.5;
    }
    .ok {
      color: var(--ok);
      font-weight: 700;
    }
    .note {
      margin-top: 12px;
      font-size: 13px;
    }
  </style>
</head>
<body>
  <main class="card">
    <span class="badge">Installation completed</span>
    <h1><span class="ok">Success.</span> Your installation is complete.</h1>
    <p>The AppMax integration token was confirmed successfully.</p>
    <p>You can close this browser tab now.</p>
    <p class="note">This page can be safely closed after installation confirmation.</p>
  </main>
</body>
</html>`

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
	accept := strings.ToLower(ctx.Request().Header("Accept"))
	wantsHTML := strings.Contains(accept, "text/html")

	token := strings.TrimSpace(ctx.Request().Query("token", ""))
	if token == "" {
		if wantsHTML {
			return NewHealthController().CallbackFrontend(ctx)
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
	if wantsHTML {
		return c.InstallCompletedFrontend(ctx)
	}
	return ctx.Response().Json(200, responses.InstallCallbackResponse{ExternalID: inst.ExternalID})
}

func (c *InstallController) InstallCompletedFrontend(ctx http.Context) http.Response {
	return ctx.Response().Data(200, "text/html; charset=utf-8", []byte(installCompletedFrontendTemplate))
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

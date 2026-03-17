package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/controllers"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/requests"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/http/responses"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/models"
	"github.com/geovannegallinati/AppStore-Appmax-App-Integration/app/services"
	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAdminURL   = "https://admin.appmax.com.br"
	testAppURL     = "https://app.example.com"
	testAppUUID    = "test-app-uuid"
	testAppNumeric = "42"
)

type mockAppmaxSvcInstall struct {
	noopAppmaxService
	authorizeFunc             func(context.Context, string, string, string) (string, error)
	generateMerchantCredsFunc func(context.Context, string) (string, string, error)
}

func (m *mockAppmaxSvcInstall) Authorize(ctx context.Context, appID, externalKey, callbackURL string) (string, error) {
	if m.authorizeFunc != nil {
		return m.authorizeFunc(ctx, appID, externalKey, callbackURL)
	}
	return "", nil
}

func (m *mockAppmaxSvcInstall) GenerateMerchantCreds(ctx context.Context, token string) (string, string, error) {
	if m.generateMerchantCredsFunc != nil {
		return m.generateMerchantCredsFunc(ctx, token)
	}
	return "", "", nil
}

type mockInstallSvcCapture struct {
	upsertFunc func(context.Context, services.UpsertInstallationInput) (*models.Installation, bool, error)
	lastInput  services.UpsertInstallationInput
}

func (m *mockInstallSvcCapture) Upsert(ctx context.Context, input services.UpsertInstallationInput) (*models.Installation, bool, error) {
	m.lastInput = input
	if m.upsertFunc != nil {
		return m.upsertFunc(ctx, input)
	}
	return &models.Installation{ExternalID: "default-ext-id"}, true, nil
}

type fakeInstallState struct {
	AppID       string
	ExternalKey string
}

func newTestInstallController(t *testing.T, appmaxSvc services.AppmaxService, installSvc services.InstallService) *controllers.InstallController {
	t.Helper()
	ctrl, err := controllers.NewInstallController(appmaxSvc, installSvc, testAdminURL, testAppURL, testAppUUID, testAppNumeric)
	require.NoError(t, err)
	return ctrl
}

func assertViewData(t *testing.T, captured capturedResponse) map[string]any {
	t.Helper()
	b, err := json.Marshal(captured.viewData)
	require.NoError(t, err)
	var data map[string]any
	require.NoError(t, json.Unmarshal(b, &data))
	return data
}

func TestInstallController_Start_MissingParams_HTML(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		queryParams: map[string]string{},
		headers:     map[string]string{"Accept": "text/html"},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Start(ctx)

	require.Equal(t, "view", ctx.resp.captured.kind)
	assert.Equal(t, "frontend/page.tmpl", ctx.resp.captured.viewTemplate)
	data := assertViewData(t, ctx.resp.captured)
	assert.Equal(t, true, data["ShowInstallForm"])
	assert.Equal(t, testAppUUID, data["DefaultInstallAppID"])
	assert.Equal(t, "", data["DefaultExternalKey"])
}

func TestInstallController_Start_MissingAppID_API(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		queryParams: map[string]string{"external_key": "some-key"},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Start(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 400, ctx.resp.captured.status)
	body, ok := ctx.resp.captured.jsonBody.(responses.MessageResponse)
	require.True(t, ok)
	assert.Equal(t, "app_id is required", body.Message)
}

func TestInstallController_Start_MissingExternalKey_API(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		queryParams: map[string]string{"app_id": testAppUUID},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Start(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 400, ctx.resp.captured.status)
	body, ok := ctx.resp.captured.jsonBody.(responses.MessageResponse)
	require.True(t, ok)
	assert.Equal(t, "external_key is required", body.Message)
}

func TestInstallController_Start_HappyPath(t *testing.T) {
	const returnedHash = "start-happy-hash"

	var capturedExternalKey string
	appmaxSvc := &mockAppmaxSvcInstall{
		authorizeFunc: func(_ context.Context, _, externalKey, _ string) (string, error) {
			capturedExternalKey = externalKey
			return returnedHash, nil
		},
	}
	ctrl := newTestInstallController(t, appmaxSvc, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		queryParams: map[string]string{
			"app_id":       testAppUUID,
			"external_key": "1234-some-uuid",
		},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Start(ctx)

	assert.Equal(t, "1234-some-uuid", capturedExternalKey)
	require.Equal(t, "redirect", ctx.resp.captured.kind)
	assert.Equal(t, 302, ctx.resp.captured.status)
	assert.Contains(t, ctx.resp.captured.redirectURL, returnedHash)

	cached := facades.Cache().GetString("install:" + returnedHash)
	require.NotEmpty(t, cached)
	var state fakeInstallState
	require.NoError(t, json.Unmarshal([]byte(cached), &state))
	assert.Equal(t, testAppUUID, state.AppID)
	assert.Equal(t, "1234-some-uuid", state.ExternalKey)
}

func TestInstallController_Start_AuthorizeFails(t *testing.T) {
	appmaxSvc := &mockAppmaxSvcInstall{
		authorizeFunc: func(context.Context, string, string, string) (string, error) {
			return "", errors.New("upstream timeout")
		},
	}
	ctrl := newTestInstallController(t, appmaxSvc, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		queryParams: map[string]string{
			"app_id":       testAppUUID,
			"external_key": "some-key",
		},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Start(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 502, ctx.resp.captured.status)
}

func TestInstallController_CallbackGuide_NoToken_HTML(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		queryParams: map[string]string{},
		headers:     map[string]string{"Accept": "text/html"},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.CallbackGuide(ctx)

	require.Equal(t, "view", ctx.resp.captured.kind)
	assert.Equal(t, "frontend/page.tmpl", ctx.resp.captured.viewTemplate)
}

func TestInstallController_CallbackGuide_NoToken_API(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{queryParams: map[string]string{}}
	ctx := newFakeHTTPContext(req)

	ctrl.CallbackGuide(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 400, ctx.resp.captured.status)
}

func TestInstallController_CallbackGuide_NoCachedState_HTML(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		queryParams: map[string]string{"token": "no-cache-token-xyz"},
		headers:     map[string]string{"Accept": "text/html"},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.CallbackGuide(ctx)

	require.Equal(t, "view", ctx.resp.captured.kind)
	assert.Equal(t, "frontend/install_completed.tmpl", ctx.resp.captured.viewTemplate)
}

func TestInstallController_CallbackGuide_HappyPath(t *testing.T) {
	const token = "guide-happy-path-token"

	state := fakeInstallState{AppID: testAppUUID, ExternalKey: "1234-test-uuid"}
	stateJSON, _ := json.Marshal(state)
	require.NoError(t, facades.Cache().Put("install:"+token, string(stateJSON), time.Minute))

	installSvc := &mockInstallSvcCapture{}
	appmaxSvc := &mockAppmaxSvcInstall{
		generateMerchantCredsFunc: func(_ context.Context, tok string) (string, string, error) {
			assert.Equal(t, token, tok)
			return "mc-client-id", "mc-client-secret", nil
		},
	}
	ctrl := newTestInstallController(t, appmaxSvc, installSvc)

	req := &fakeHTTPRequest{
		queryParams: map[string]string{"token": token},
		headers:     map[string]string{"Accept": "text/html"},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.CallbackGuide(ctx)

	require.Equal(t, "view", ctx.resp.captured.kind)
	assert.Equal(t, "frontend/install_completed.tmpl", ctx.resp.captured.viewTemplate)
	assert.Equal(t, testAppUUID, installSvc.lastInput.AppID)
	assert.Equal(t, "1234-test-uuid", installSvc.lastInput.ExternalKey)
	assert.Equal(t, "mc-client-id", installSvc.lastInput.MerchantClientID)
	assert.Equal(t, "mc-client-secret", installSvc.lastInput.MerchantClientSecret)
	assert.Empty(t, facades.Cache().GetString("install:"+token), "cache should be cleared after use")
}

func TestInstallController_CallbackGuide_AppIDMismatch(t *testing.T) {
	const token = "guide-mismatch-token"

	state := fakeInstallState{AppID: "wrong-app-id", ExternalKey: "some-key"}
	stateJSON, _ := json.Marshal(state)
	require.NoError(t, facades.Cache().Put("install:"+token, string(stateJSON), time.Minute))

	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{queryParams: map[string]string{"token": token}}
	ctx := newFakeHTTPContext(req)

	ctrl.CallbackGuide(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 400, ctx.resp.captured.status)
	body, ok := ctx.resp.captured.jsonBody.(responses.MessageResponse)
	require.True(t, ok)
	assert.Equal(t, "invalid app_id", body.Message)
}

func TestInstallController_CallbackGuide_GenerateCredsFails(t *testing.T) {
	const token = "guide-creds-fail-token"

	state := fakeInstallState{AppID: testAppUUID, ExternalKey: "some-key"}
	stateJSON, _ := json.Marshal(state)
	require.NoError(t, facades.Cache().Put("install:"+token, string(stateJSON), time.Minute))

	appmaxSvc := &mockAppmaxSvcInstall{
		generateMerchantCredsFunc: func(context.Context, string) (string, string, error) {
			return "", "", errors.New("gateway timeout")
		},
	}
	ctrl := newTestInstallController(t, appmaxSvc, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{queryParams: map[string]string{"token": token}}
	ctx := newFakeHTTPContext(req)

	ctrl.CallbackGuide(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 502, ctx.resp.captured.status)
}

func TestInstallController_Callback_HappyPath(t *testing.T) {
	installSvc := &mockInstallSvcCapture{
		upsertFunc: func(_ context.Context, _ services.UpsertInstallationInput) (*models.Installation, bool, error) {
			return &models.Installation{ExternalID: "ext-abc"}, true, nil
		},
	}
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, installSvc)

	req := &fakeHTTPRequest{
		bindResult: requests.InstallCallbackRequest{
			AppID:                testAppNumeric,
			ExternalKey:          "some-key",
			MerchantClientID:     "mc-id",
			MerchantClientSecret: "mc-secret",
		},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Callback(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 200, ctx.resp.captured.status)
	body, ok := ctx.resp.captured.jsonBody.(responses.InstallCallbackResponse)
	require.True(t, ok)
	assert.Equal(t, "ext-abc", body.ExternalID)
	assert.Equal(t, testAppNumeric, installSvc.lastInput.AppID)
	assert.Equal(t, "some-key", installSvc.lastInput.ExternalKey)
	assert.Equal(t, "mc-id", installSvc.lastInput.MerchantClientID)
	assert.Equal(t, "mc-secret", installSvc.lastInput.MerchantClientSecret)
}

func TestInstallController_Callback_MissingFields(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		bindResult: requests.InstallCallbackRequest{
			AppID:       testAppNumeric,
			ExternalKey: "some-key",
		},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Callback(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 400, ctx.resp.captured.status)
}

func TestInstallController_Callback_WrongAppID(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{
		bindResult: requests.InstallCallbackRequest{
			AppID:                "wrong-app-id",
			ExternalKey:          "some-key",
			MerchantClientID:     "mc-id",
			MerchantClientSecret: "mc-secret",
		},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Callback(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 400, ctx.resp.captured.status)
	body, ok := ctx.resp.captured.jsonBody.(responses.MessageResponse)
	require.True(t, ok)
	assert.Equal(t, "invalid app_id", body.Message)
}

func TestInstallController_Callback_UpsertFails(t *testing.T) {
	installSvc := &mockInstallSvcCapture{
		upsertFunc: func(context.Context, services.UpsertInstallationInput) (*models.Installation, bool, error) {
			return nil, false, errors.New("db error")
		},
	}
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, installSvc)

	req := &fakeHTTPRequest{
		bindResult: requests.InstallCallbackRequest{
			AppID:                testAppNumeric,
			ExternalKey:          "some-key",
			MerchantClientID:     "mc-id",
			MerchantClientSecret: "mc-secret",
		},
	}
	ctx := newFakeHTTPContext(req)

	ctrl.Callback(ctx)

	require.Equal(t, "json", ctx.resp.captured.kind)
	assert.Equal(t, 500, ctx.resp.captured.status)
}

func TestInstallController_InstallStartFrontend_DefaultsAppID(t *testing.T) {
	ctrl := newTestInstallController(t, &mockAppmaxSvcInstall{}, &mockInstallSvcCapture{})

	req := &fakeHTTPRequest{}
	ctx := newFakeHTTPContext(req)

	ctrl.InstallStartFrontend(ctx, "", "")

	require.Equal(t, "view", ctx.resp.captured.kind)
	assert.Equal(t, "frontend/page.tmpl", ctx.resp.captured.viewTemplate)
	data := assertViewData(t, ctx.resp.captured)
	assert.Equal(t, testAppUUID, data["DefaultInstallAppID"])
	assert.Equal(t, "", data["DefaultExternalKey"])
	assert.Equal(t, true, data["ShowInstallForm"])
}

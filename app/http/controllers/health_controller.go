package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
)

type HealthController struct {
	appURL string
}

func NewHealthController(appURL string) *HealthController {
	return &HealthController{
		appURL: appURL,
	}
}

func (c *HealthController) Check(ctx contractshttp.Context) contractshttp.Response {
	if requestWantsHTML(ctx) {
		return c.HealthFrontend(ctx)
	}

	return ctx.Response().Json(200, contractshttp.Json{"status": "ok"})
}

func (c *HealthController) FrontendCheck(ctx contractshttp.Context) contractshttp.Response {
	switch ctx.Request().Path() {
	case "/health":
		return c.HealthFrontend(ctx)
	case "/integrations/appmax/callback/install":
		return c.CallbackFrontend(ctx)
	default:
		return c.RootFrontend(ctx)
	}
}

func (c *HealthController) RootFrontend(ctx contractshttp.Context) contractshttp.Response {
	baseURL := frontendBaseURL(ctx, c.appURL)

	page := frontendPageData{
		Title:       "AppMax Checkout Demo - Frontend",
		Badge:       "Frontend",
		Headline:    "Application is online and reachable",
		Message:     "The app public URL is responding correctly.",
		Submessage:  "Use the shortcuts below to validate install, callback, and webhook.",
		ActiveRoute: routeRoot,
		PageKind:    "root",
		Endpoints:   frontendEndpoints(baseURL, routeRoot),
		Buttons: []frontendAction{
			{Label: "Start install", URL: routeInstallStart},
			{Label: "Configure webhook", URL: routeWebhook, Secondary: true},
		},
		Tips: []string{
			"Open /health in a browser to view the health frontend page.",
			"The /health endpoint still returns JSON for API calls without Accept text/html.",
		},
	}
	return c.renderFrontend(ctx, page)
}

func (c *HealthController) HealthFrontend(ctx contractshttp.Context) contractshttp.Response {
	baseURL := frontendBaseURL(ctx, c.appURL)

	page := frontendPageData{
		Title:       "AppMax Checkout Demo - Health",
		Badge:       "Health",
		Headline:    "Health endpoint validated",
		Message:     "This page confirms the health frontend is active.",
		Submessage:  "For API probes, call /health without Accept text/html to receive JSON.",
		ActiveRoute: routeHealth,
		PageKind:    "health",
		Endpoints:   frontendEndpoints(baseURL, routeHealth),
		Tips: []string{
			"Expected health API response: {\"status\":\"ok\"}.",
		},
	}
	return c.renderFrontend(ctx, page)
}

func (c *HealthController) CallbackFrontend(ctx contractshttp.Context) contractshttp.Response {
	baseURL := frontendBaseURL(ctx, c.appURL)

	page := frontendPageData{
		Title:       "AppMax Checkout Demo - Callback",
		Badge:       "Callback",
		Headline:    "Install callback endpoint is ready",
		Message:     "The callback URL is reachable and ready for the Appmax flow.",
		Submessage:  "The same endpoint also receives install confirmation POST requests.",
		ActiveRoute: routeInstallCallback,
		PageKind:    "callback",
		Endpoints:   frontendEndpoints(baseURL, routeInstallCallback),
		Tips: []string{
			"When a token is present in the URL, the callback completes installation automatically.",
			"Without a token and with a browser, this page is shown for manual validation.",
		},
	}
	return c.renderFrontend(ctx, page)
}

func (c *HealthController) renderFrontend(ctx contractshttp.Context, page frontendPageData) contractshttp.Response {
	return ctx.Response().View().Make("frontend/page.tmpl", page)
}

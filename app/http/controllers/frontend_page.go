package controllers

import (
	"fmt"
	neturl "net/url"
	"strings"

	contractshttp "github.com/goravel/framework/contracts/http"
)

const (
	routeRoot            = "/"
	routeHealth          = "/health"
	routeInstallStart    = "/install/start"
	routeInstallCallback = "/integrations/appmax/callback/install"
	routeWebhook         = "/webhooks/appmax"
)

type frontendEndpoint struct {
	Label       string
	Method      string
	Path        string
	Description string
	URL         string
	Active      bool
}

type frontendAction struct {
	Label     string
	URL       string
	Secondary bool
}

type frontendPageData struct {
	Title               string
	Badge               string
	Headline            string
	Message             string
	Submessage          string
	ActiveRoute         string
	PageKind            string
	Endpoints           []frontendEndpoint
	Tips                []string
	Buttons             []frontendAction
	ShowInstallForm     bool
	InstallFormAction   string
	DefaultInstallAppID string
	DefaultExternalKey  string
	ShowWebhookGuide    bool
	AppmaxEnvironment   string
	AppmaxAdminURL      string
	AppmaxApphookURL    string
	WebhookEndpointURL  string
}

func requestWantsHTML(ctx contractshttp.Context) bool {
	return strings.Contains(strings.ToLower(ctx.Request().Header("Accept")), "text/html")
}

func frontendBaseURL(ctx contractshttp.Context, fallback string) string {
	if fromURL := parseBaseURL(ctx.Request().Url()); fromURL != "" {
		return fromURL
	}
	if fromFullURL := parseBaseURL(ctx.Request().FullUrl()); fromFullURL != "" {
		return fromFullURL
	}

	host := strings.TrimSpace(ctx.Request().Host())
	if host != "" {
		scheme := strings.TrimSpace(ctx.Request().Header("X-Forwarded-Proto"))
		if scheme == "" {
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(fallback)), "https://") {
				scheme = "https"
			} else {
				scheme = "http"
			}
		}

		return strings.TrimRight(scheme+"://"+host, "/")
	}

	return strings.TrimRight(strings.TrimSpace(fallback), "/")
}

func parseBaseURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	parsed, err := neturl.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}

	return strings.TrimRight(parsed.Scheme+"://"+parsed.Host, "/")
}

func absoluteURL(baseURL, endpointPath string) string {
	if strings.TrimSpace(baseURL) == "" {
		return endpointPath
	}

	if strings.HasPrefix(endpointPath, "/") {
		return strings.TrimRight(baseURL, "/") + endpointPath
	}

	return strings.TrimRight(baseURL, "/") + "/" + endpointPath
}

func frontendEndpoints(baseURL, activeRoute string) []frontendEndpoint {
	endpoints := []frontendEndpoint{
		{
			Label:       "Frontend",
			Method:      "GET",
			Path:        routeRoot,
			Description: "Main visual validation page for the app.",
		},
		{
			Label:       "Health",
			Method:      "GET",
			Path:        routeHealth,
			Description: "Health frontend page for availability validation.",
		},
		{
			Label:       "Install",
			Method:      "GET",
			Path:        routeInstallStart,
			Description: "Install page with a button to start the flow.",
		},
		{
			Label:       "Callback",
			Method:      "GET/POST",
			Path:        routeInstallCallback,
			Description: "Install callback return and Appmax confirmation endpoint.",
		},
		{
			Label:       "Webhook",
			Method:      "GET/POST",
			Path:        routeWebhook,
			Description: "Webhook setup frontend and Appmax POST events endpoint.",
		},
	}

	for i := range endpoints {
		endpoints[i].Active = endpoints[i].Path == activeRoute
		endpoints[i].URL = absoluteURL(baseURL, endpoints[i].Path)
	}

	return endpoints
}

func appmaxEnvironmentName(adminURL string) string {
	if strings.Contains(strings.ToLower(adminURL), "sandbox") {
		return "Sandbox"
	}

	return "Production"
}

func apphookURLFromAdmin(adminURL string) string {
	normalized := strings.TrimRight(strings.TrimSpace(adminURL), "/")
	if normalized == "" {
		return ""
	}

	return fmt.Sprintf("%s/v2/apphook", normalized)
}

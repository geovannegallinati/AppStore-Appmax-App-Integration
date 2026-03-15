package controllers

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"strings"
)

type HealthController struct{}

type frontendPage struct {
	Title       string
	Badge       string
	Headline    string
	Message     string
	Submessage  string
	Tip         string
	ActiveRoute string
}

const frontendTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>__TITLE__</title>
  <style>
    :root {
      --bg: #0d1b2a;
      --panel: #1b263b;
      --ok: #2ec4b6;
      --txt: #e0e1dd;
      --muted: #9aa6b2;
      --accent: #ffb703;
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
    }
    .card {
      width: min(92vw, 760px);
      border: 1px solid #34445e;
      background: linear-gradient(160deg, rgba(27,38,59,0.95), rgba(17,28,43,0.95));
      border-radius: 16px;
      padding: 28px 24px;
      box-shadow: 0 14px 40px rgba(0,0,0,0.35);
    }
    .pill {
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
      font-size: clamp(24px, 3.2vw, 36px);
      line-height: 1.15;
    }
    p {
      margin: 0 0 12px;
      color: var(--muted);
      line-height: 1.5;
    }
    .ok { color: var(--ok); font-weight: 700; }
    .tip {
      margin-top: 18px;
      padding: 12px 14px;
      border-left: 4px solid var(--accent);
      background: rgba(255,183,3,0.08);
      color: #f1d28a;
      border-radius: 8px;
    }
    .endpoints {
      margin-top: 18px;
      border: 1px solid #32455f;
      background: rgba(255,255,255,0.03);
      border-radius: 12px;
      padding: 14px;
    }
    .endpoints h2 {
      margin: 0 0 10px;
      font-size: 16px;
      color: #d4deea;
    }
    .endpoint {
      margin: 0 0 10px;
      padding-bottom: 10px;
      border-bottom: 1px solid rgba(255,255,255,0.08);
    }
    .endpoint:last-child {
      margin-bottom: 0;
      padding-bottom: 0;
      border-bottom: 0;
    }
    .endpoint.active {
      border-left: 3px solid var(--ok);
      padding-left: 10px;
    }
    .endpoint a {
      color: #9cd6ff;
      text-decoration: none;
      font-weight: 700;
      word-break: break-all;
    }
    .endpoint a:hover { text-decoration: underline; }
    .endpoint p {
      margin: 6px 0 0;
      font-size: 14px;
    }
    code {
      color: #fff;
      background: rgba(255,255,255,0.08);
      padding: 2px 6px;
      border-radius: 6px;
    }
  </style>
</head>
<body>
  <main class="card">
    <span class="pill">__BADGE__</span>
    <h1>__HEADLINE__ <span class="ok">working</span></h1>
    <p>__MESSAGE__</p>
    <p>__SUBMESSAGE__</p>
    <section class="endpoints">
      <h2>Available Frontend Endpoints</h2>
      __ENDPOINTS__
    </section>
    <div class="tip">__TIP__</div>
  </main>
</body>
</html>`

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Check(ctx contractshttp.Context) contractshttp.Response {
	accept := strings.ToLower(ctx.Request().Header("Accept"))
	if strings.Contains(accept, "text/html") {
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
	page := frontendPage{
		Title:       "AppMax Checkout Demo - Main Frontend",
		Badge:       "Main Frontend",
		Headline:    "Root endpoint is",
		Message:     "You reached the main frontend endpoint through your app + nginx + ngrok stack.",
		Submessage:  "This page confirms the public root URL is online and serving HTML.",
		Tip:         `Tip: open <code>/health</code> in browser for the health frontend page.`,
		ActiveRoute: "/",
	}
	return c.renderFrontend(ctx, page)
}

func (c *HealthController) HealthFrontend(ctx contractshttp.Context) contractshttp.Response {
	page := frontendPage{
		Title:       "AppMax Checkout Demo - Health Frontend",
		Badge:       "Health Frontend",
		Headline:    "Health frontend endpoint is",
		Message:     "You reached the /health frontend page.",
		Submessage:  "For API probes, call /health with a non-HTML Accept header to receive JSON status.",
		Tip:         `Tip: this browser page is visual; API health checks should expect <code>{"status":"ok"}</code>.`,
		ActiveRoute: "/health",
	}
	return c.renderFrontend(ctx, page)
}

func (c *HealthController) CallbackFrontend(ctx contractshttp.Context) contractshttp.Response {
	page := frontendPage{
		Title:       "AppMax Checkout Demo - Callback Frontend",
		Badge:       "Callback Frontend",
		Headline:    "Install callback endpoint is",
		Message:     "You reached the installation callback endpoint frontend.",
		Submessage:  "This page confirms callback URL reachability over the active public tunnel.",
		Tip:         `Tip: this is the callback endpoint at <code>/integrations/appmax/callback/install</code>.`,
		ActiveRoute: "/integrations/appmax/callback/install",
	}
	return c.renderFrontend(ctx, page)
}

func (c *HealthController) renderFrontend(ctx contractshttp.Context, page frontendPage) contractshttp.Response {
	html := strings.NewReplacer(
		"__TITLE__", page.Title,
		"__BADGE__", page.Badge,
		"__HEADLINE__", page.Headline,
		"__MESSAGE__", page.Message,
		"__SUBMESSAGE__", page.Submessage,
		"__TIP__", page.Tip,
		"__ENDPOINTS__", c.endpointsHTML(page.ActiveRoute),
	).Replace(frontendTemplate)

	return ctx.Response().Data(200, "text/html; charset=utf-8", []byte(html))
}

func (c *HealthController) endpointsHTML(activeRoute string) string {
	return endpointHTML("/", "This is the main frontend endpoint.", activeRoute) +
		endpointHTML("/health", "This is the health frontend endpoint.", activeRoute) +
		endpointHTML("/integrations/appmax/callback/install", "This is the callback endpoint.", activeRoute)
}

func endpointHTML(path string, description string, activeRoute string) string {
	className := "endpoint"
	if path == activeRoute {
		className += " active"
	}

	return `<div class="` + className + `"><a href="` + path + `">` + path + `</a><p>` + description + `</p></div>`
}

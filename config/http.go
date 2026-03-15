package config

import (
	"github.com/goravel/framework/facades"

	routecontract "github.com/goravel/framework/contracts/route"
	ginfacades    "github.com/goravel/gin/facades"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		"default":         "gin",
		"host":            config.Env("APP_HOST", "0.0.0.0"),
		"port":            config.Env("APP_PORT", "8080"),
		"request_timeout": 120,
		"drivers": map[string]any{
			"gin": map[string]any{
				"body_limit":   4096,
				"header_limit": 4096,
				"route": func() (routecontract.Route, error) {
					return ginfacades.Route("gin"), nil
				},
			},
		},
		"tls": map[string]any{
			"host": config.Env("APP_TLS_HOST", ""),
			"port": config.Env("APP_TLS_PORT", ""),
			"ssl": map[string]any{
				"cert": config.Env("APP_TLS_CERT", ""),
				"key":  config.Env("APP_TLS_KEY", ""),
			},
		},
	})
}

package config

import "github.com/goravel/framework/facades"

func init() {
	config := facades.Config()
	config.Add("app", map[string]any{
		"name":     config.Env("APP_NAME", "AppMax Checkout"),
		"env":      config.Env("APP_ENV", "production"),
		"debug":    config.Env("APP_DEBUG", false),
		"url":      config.Env("APP_URL", "http://localhost:8080"),
		"key":      config.Env("APP_KEY", ""),
		"timezone": config.Env("APP_TIMEZONE", "UTC"),
		"locale":   "pt_BR",
	})
}

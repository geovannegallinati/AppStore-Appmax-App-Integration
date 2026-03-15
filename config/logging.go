package config

import "github.com/goravel/framework/facades"

func init() {
	config := facades.Config()
	config.Add("logging", map[string]any{
		"default": config.Env("LOG_CHANNEL", "stack"),
		"channels": map[string]any{
			"stack": map[string]any{
				"driver":   "stack",
				"channels": []string{"single"},
			},
			"single": map[string]any{
				"driver": "single",
				"path":   config.Env("LOG_PATH", "storage/logs/goravel.log"),
				"level":  config.Env("LOG_LEVEL", "debug"),
				"days":   14,
				"print":  config.Env("LOG_PRINT", true),
			},
		},
	})
}

package config

import (
	"github.com/goravel/framework/facades"

	drivercontract "github.com/goravel/framework/contracts/database/driver"
	postgresfacades "github.com/goravel/postgres/facades"
)

func init() {
	config := facades.Config()
	config.Add("database", map[string]any{
		"default": config.Env("DB_CONNECTION", "postgres"),
		"connections": map[string]any{
			"postgres": map[string]any{
				"driver":   "postgres",
				"host":     config.Env("DB_HOST", "localhost"),
				"port":     config.Env("DB_PORT", 5432),
				"database": config.Env("DB_DATABASE", "appmax_checkout"),
				"username": config.Env("DB_USERNAME", "appmax"),
				"password": config.Env("DB_PASSWORD", ""),
				"sslmode":  config.Env("DB_SSLMODE", "disable"),
				"schema":   config.Env("DB_SCHEMA", "public"),
				"timezone": config.Env("DB_TIMEZONE", "UTC"),
				"prefix":   "",
				"singular": false,
				"via": func() (drivercontract.Driver, error) {
					return postgresfacades.Postgres("postgres")
				},
			},
		},
		"migrations": map[string]any{
			"table": "migrations",
		},
		"pool": map[string]any{
			"max_idle_conns":    5,
			"max_open_conns":    25,
			"conn_max_lifetime": 300,
		},
		"redis": map[string]any{
			"default": map[string]any{
				"host":     config.Env("REDIS_HOST", "localhost"),
				"port":     config.Env("REDIS_PORT", "6379"),
				"username": config.Env("REDIS_USERNAME", ""),
				"password": config.Env("REDIS_PASSWORD", ""),
				"database": config.Env("REDIS_DB", 0),
			},
		},
	})
}

package services

import (
	"context"
	"time"
)

// Logger abstracts logging for application services.
// Framework-specific logging should stay in adapters/bootstrap.
type Logger interface {
	Warningf(format string, args ...any)
	Errorf(format string, args ...any)
}

// MerchantTokenCache abstracts cache operations used by token manager.
type MerchantTokenCache interface {
	GetString(ctx context.Context, key string) string
	Put(ctx context.Context, key, value string, ttl time.Duration) error
}

type noopLogger struct{}

func (noopLogger) Warningf(string, ...any) {}

func (noopLogger) Errorf(string, ...any) {}

type noopMerchantTokenCache struct{}

func (noopMerchantTokenCache) GetString(context.Context, string) string {
	return ""
}

func (noopMerchantTokenCache) Put(context.Context, string, string, time.Duration) error {
	return nil
}

func coalesceLogger(logger Logger) Logger {
	if logger == nil {
		return noopLogger{}
	}

	return logger
}

func coalesceMerchantTokenCache(cache MerchantTokenCache) MerchantTokenCache {
	if cache == nil {
		return noopMerchantTokenCache{}
	}

	return cache
}

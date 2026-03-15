package goravel

import (
	"context"
	"time"

	"github.com/goravel/framework/facades"
)

type Cache struct{}

func NewCache() *Cache {
	return &Cache{}
}

func (Cache) GetString(ctx context.Context, key string) string {
	return facades.Cache().WithContext(ctx).GetString(key)
}

func (Cache) Put(ctx context.Context, key, value string, ttl time.Duration) error {
	return facades.Cache().WithContext(ctx).Put(key, value, ttl)
}

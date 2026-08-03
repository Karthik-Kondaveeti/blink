package cache

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("CacheMiss")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Close() error
}

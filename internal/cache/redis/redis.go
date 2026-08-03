package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Karthik-Kondaveeti/blink/internal/cache"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(address string, password string, db int) (*Redis, error) {
	if address == "" {
		return nil, errors.New("redis address is empty")
	}

	client := redis.NewClient(
		&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		},
	)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return &Redis{
		client: client,
	}, nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", cache.ErrCacheMiss
		}
		return "", fmt.Errorf("Error getting key: %w", err)
	}
	return value, nil
}

func (r *Redis) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return errors.New("Error setting key: " + err.Error())
	}
	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

package database

import (
	"context"
)

type Database interface {
	AddLink(ctx context.Context, id uint64, originalURL string) error
	GetLink(ctx context.Context, id uint64) (string, error)
	Close() error
}

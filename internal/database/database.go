package database

import (
	"context"
)

type Database interface {
	AddLink(ctx context.Context, shortCode string, originalURL string) error
	GetLink(ctx context.Context, shortCode string) (string, error)
	Close() error
}

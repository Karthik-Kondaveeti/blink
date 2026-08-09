package service

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"time"

	"github.com/Karthik-Kondaveeti/blink/internal/cache"
	"github.com/Karthik-Kondaveeti/blink/internal/database"
	"github.com/Karthik-Kondaveeti/blink/internal/generator"
)

type Service struct {
	db        database.Database
	cdb       cache.Cache
	generator generator.Generator
	logger    *slog.Logger
}

func New(db database.Database, cdb cache.Cache, generator generator.Generator, logger *slog.Logger) (*Service, error) {
	return &Service{
		db:        db,
		cdb:       cdb,
		generator: generator,
		logger:    logger,
	}, nil
}

func isValidURL(link string) bool {
	u, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}

func (s *Service) GetLink(ctx context.Context, shortCode string) (string, error) {
	cacheKey := "link:" + shortCode
	originalURL, err := s.cdb.Get(ctx, cacheKey)

	// Cache Hit
	if err == nil {
		s.logger.Debug("cache hit", "short_code", shortCode)
		return originalURL, nil
	}

	if errors.Is(err, cache.ErrCacheMiss) {
		// Cache Miss
		s.logger.Debug(
			"cache miss",
			"short_code", shortCode,
		)

		id, err := s.generator.Decode(shortCode)
		if err != nil {
			return "", err
		}
		originalURL, err = s.db.GetLink(ctx, id)
		if err != nil {
			return "", err
		}
		if err := s.cdb.Set(ctx, cacheKey, originalURL, time.Hour); err != nil {
			s.logger.Warn(
				"failed to cache link",
				"short_code", shortCode,
				"error", err,
			)
		}
		return originalURL, nil
	} else {
		s.logger.Error("redis get failed", "error", err)
		return "", err
	}
}

func (s *Service) AddLink(ctx context.Context, originalURL string) (string, error) {
	if !isValidURL(originalURL) {
		return "", errors.New("Invalid URL!")
	}

	id := s.generator.Generate()
	err := s.db.AddLink(
		ctx,
		id,
		originalURL,
	)

	if err != nil {
		return "", err
	}
	shortCode := s.generator.Encode(id)
	return shortCode, nil
}

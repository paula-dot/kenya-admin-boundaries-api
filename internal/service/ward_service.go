package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
)

// wardRepo defines the repository subset required by the WardService
type wardRepo interface {
	ListWards(ctx context.Context, limit, offset int32) ([]postgres.ListWardsRow, error)
}

// WardService handles business logic related to Wards
type WardService struct {
	repo  wardRepo
	cache CacheRepository
}

// NewWardService constructs a new WardService
func NewWardService(repo wardRepo, cache CacheRepository) *WardService {
	return &WardService{
		repo:  repo,
		cache: cache,
	}
}

// ListWards fetches a paginated list of Wards (no geometry)
func (s *WardService) ListWards(ctx context.Context, limit, offset int32) ([]postgres.ListWardsRow, error) {
	cacheKey := fmt.Sprintf("wards:list:page:%d:%d", limit, offset)

	if s.cache != nil {
		if cachedData, err := s.cache.Get(ctx, cacheKey); err == nil && cachedData != nil {
			var units []postgres.ListWardsRow
			if err := json.Unmarshal(cachedData, &units); err == nil {
				return units, nil
			}
		}
	}

	rows, err := s.repo.ListWards(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list wards: %w", err)
	}

	if s.cache != nil {
		if collBytes, err := json.Marshal(rows); err == nil {
			go func(ctx context.Context, key string, b []byte) {
				_ = s.cache.Set(ctx, key, b, 24*time.Hour)
			}(context.Background(), cacheKey, collBytes)
		}
	}

	return rows, nil
}

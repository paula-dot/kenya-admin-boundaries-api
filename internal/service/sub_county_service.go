package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
)

// subCountyRepo defines the methods we need from the database package.
type subCountyRepo interface {
	GetAllSubCounties(ctx context.Context) ([]*domain.SubCounty, error)
	GetSubCountiesByCounty(ctx context.Context, countyCode string) ([]*domain.SubCounty, error)
}

// SubCountyService handles business logic related to Sub-Counties.
type SubCountyService struct {
	repo  subCountyRepo
	cache CacheRepository
}

// NewSubCountyService initializes a new sub-county service.
func NewSubCountyService(repo subCountyRepo, cache CacheRepository) *SubCountyService {
	return &SubCountyService{
		repo:  repo,
		cache: cache,
	}
}

// GetAll returns a complete list of all sub-counties.
func (s *SubCountyService) GetAll(ctx context.Context) ([]*domain.SubCounty, error) {
	cacheKey := "subcounties:all"

	if s.cache != nil {
		cachedData, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			var cachedSubCounties []*domain.SubCounty
			if err := json.Unmarshal(cachedData, &cachedSubCounties); err == nil {
				return cachedSubCounties, nil
			}
		}
	}

	subCounties, err := s.repo.GetAllSubCounties(ctx)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if outBytes, err := json.Marshal(subCounties); err == nil {
			go func(ctx context.Context, key string, b []byte) {
				_ = s.cache.Set(ctx, key, b, 24*time.Hour)
			}(context.Background(), cacheKey, outBytes)
		}
	}

	return subCounties, nil
}

// GetByCountyCode returns sub-counties filtered by a given county code.
func (s *SubCountyService) GetByCountyCode(ctx context.Context, code string) ([]*domain.SubCounty, error) {
	cacheKey := fmt.Sprintf("subcounties:county:%s", code)

	if s.cache != nil {
		cachedData, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			var cachedSubCounties []*domain.SubCounty
			if err := json.Unmarshal(cachedData, &cachedSubCounties); err == nil {
				return cachedSubCounties, nil
			}
		}
	}

	subCounties, err := s.repo.GetSubCountiesByCounty(ctx, code)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if outBytes, err := json.Marshal(subCounties); err == nil {
			go func(ctx context.Context, key string, b []byte) {
				_ = s.cache.Set(ctx, key, b, 24*time.Hour)
			}(context.Background(), cacheKey, outBytes)
		}
	}

	return subCounties, nil
}

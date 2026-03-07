package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
	"github.com/paula-dot/kenya-admin-boundaries-api/pkg/geojson"
)

// ConstituencyListItem is the shape expected by handlers/router for constituency lists.
type ConstituencyListItem struct {
	ID       int32
	Slug     string
	Name     string
	Geometry []byte
}

// constituencyRepo defines the subset of repository methods used by the service.
type constituencyRepo interface {
	ListConstituencies(ctx context.Context) ([]*domain.Constituency, error)
	ListConstituenciesByCounty(ctx context.Context, countyCode string) ([]postgres.ListConstituenciesByCountyRow, error)
	ListConstituenciesMetadataByCounty(ctx context.Context, countyCode string) ([]postgres.ListConstituenciesMetadataByCountyRow, error)
}

// ConstituencyService provides helper methods around constituencies.
type ConstituencyService struct {
	repo  constituencyRepo
	cache CacheRepository
}

// NewConstituencyService constructs a new service.
func NewConstituencyService(repo constituencyRepo, cache CacheRepository) *ConstituencyService {
	return &ConstituencyService{
		repo:  repo,
		cache: cache,
	}
}

// ListConstituenciesByCountySlug returns constituencies for a given county code/slug.
func (s *ConstituencyService) ListConstituenciesByCountySlug(ctx context.Context, slug string) ([]struct {
	ID       int32
	Slug     string
	Name     string
	Geometry []byte
}, error) {

	// Define anonymous return structure mapped by this func
	type outT []struct {
		ID       int32
		Slug     string
		Name     string
		Geometry []byte
	}

	cacheKey := fmt.Sprintf("constituencies:features:county:%s", slug)

	// 1. Check Cache
	if s.cache != nil {
		cachedData, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			var cachedOut outT
			if err := json.Unmarshal(cachedData, &cachedOut); err == nil {
				return cachedOut, nil
			}
		}
	}

	// 2. Fetch from Repo
	rows, err := s.repo.ListConstituenciesByCounty(ctx, slug)
	if err != nil {
		return nil, err
	}

	out := make([]struct {
		ID       int32
		Slug     string
		Name     string
		Geometry []byte
	}, 0, len(rows))

	for _, r := range rows {
		var id int32 = 0
		if r.ID != "" {
			if parsed, err := strconv.Atoi(r.ID); err == nil {
				id = int32(parsed)
			}
		}
		item := struct {
			ID       int32
			Slug     string
			Name     string
			Geometry []byte
		}{
			ID:       id,
			Slug:     r.ID,
			Name:     r.Name,
			Geometry: r.Geojson,
		}
		out = append(out, item)
	}

	// 3. Populate cache asynchronously
	if s.cache != nil {
		if outBytes, err := json.Marshal(out); err == nil {
			go func(ctx context.Context, key string, b []byte) {
				_ = s.cache.Set(ctx, key, b, 24*time.Hour)
			}(context.Background(), cacheKey, outBytes)
		}
	}

	return out, nil
}

// ListConstituenciesMetadataByCountySlug returns lightweight constituencies (code and name) for a given county slug.
func (s *ConstituencyService) ListConstituenciesMetadataByCountySlug(ctx context.Context, slug string) ([]postgres.ListConstituenciesMetadataByCountyRow, error) {
	cacheKey := fmt.Sprintf("constituencies:meta:county:%s", slug)

	if s.cache != nil {
		cachedData, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			var cachedMeta []postgres.ListConstituenciesMetadataByCountyRow
			if err := json.Unmarshal(cachedData, &cachedMeta); err == nil {
				return cachedMeta, nil
			}
		}
	}

	meta, err := s.repo.ListConstituenciesMetadataByCounty(ctx, slug)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if metaBytes, err := json.Marshal(meta); err == nil {
			go func(ctx context.Context, key string, b []byte) {
				_ = s.cache.Set(ctx, key, b, 24*time.Hour)
			}(context.Background(), cacheKey, metaBytes)
		}
	}

	return meta, nil
}

// ListAllAsFeatureCollection returns all constituencies packaged as a GeoJSON FeatureCollection.
func (s *ConstituencyService) ListAllAsFeatureCollection(ctx context.Context) (*geojson.FeatureCollection, error) {
	cacheKey := "constituencies:all:fc"

	// 1. Check Cache
	if s.cache != nil {
		cachedData, err := s.cache.Get(ctx, cacheKey)
		if err == nil && cachedData != nil {
			var fc geojson.FeatureCollection
			if err := json.Unmarshal(cachedData, &fc); err == nil {
				return &fc, nil
			}
		}
	}

	// 2. Cache Miss: Fetch from DB
	constituencies, err := s.repo.ListConstituencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list constituencies: %w", err)
	}

	features := make([]geojson.Feature, 0, len(constituencies))
	for _, c := range constituencies {
		feature := geojson.Feature{
			Type: "Feature",
			ID:   fmt.Sprintf("%d", c.ID),
			Properties: map[string]interface{}{
				"code":      c.Slug,
				"name":      c.Name,
				"county_id": c.CountyID,
			},
			Geometry: json.RawMessage(c.GeoJSON),
		}
		features = append(features, feature)
	}

	collection := geojson.NewFeatureCollection(features)

	// 3. Cache asynchronously
	if s.cache != nil {
		if collBytes, err := json.Marshal(collection); err == nil {
			go func(ctx context.Context, key string, b []byte) {
				_ = s.cache.Set(ctx, key, b, 24*time.Hour)
			}(context.Background(), cacheKey, collBytes)
		}
	}

	return collection, nil
}

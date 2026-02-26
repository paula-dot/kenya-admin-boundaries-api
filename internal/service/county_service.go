package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
	"github.com/paula-dot/kenya-admin-boundaries-api/pkg/geojson"
)

// CountyRepository defines the required database operations.
// The PostgreSQL implementation will satisfy this interface.
type CountyRepository interface {
	GetCountyByID(ctx context.Context, id int32) (*domain.County, error)
	ListCounties(ctx context.Context) ([]*domain.County, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// CountyService orchestrates business logic and data formatting.
type CountyService struct {
	repo  CountyRepository
	cache CacheRepository
}

// NewCountyService is the constructor function.
func NewCountyService(repo CountyRepository, cache CacheRepository) *CountyService {
	return &CountyService{
		repo:  repo,
		cache: cache,
	}
}

// GetCountyAsFeature fetches a county and formats it as a standard GeoJSON Feature.
func (s *CountyService) GetCountyAsFeature(ctx context.Context, id int32) (*geojson.Feature, error) {
	cacheKey := fmt.Sprintf("county:feature:%d", id)

	// 1. Check the cache
	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != nil {
		// Cache hit: Unmarshal the pre-computed GeoJSON and return immediately
		var feature geojson.Feature
		if err := json.Unmarshal(cachedData, &feature); err == nil {
			return &feature, nil
		}
	}

	// 2. Cache Miss: Fetch from PostgreSQL
	county, err := s.repo.GetCountyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch county from db: %w", err)
	}

	// 3. Format into GeoJSON
	feature := &geojson.Feature{
		Type: "Feature",
		ID:   fmt.Sprintf("%d", county.ID),
		Properties: map[string]interface{}{
			"code":       county.Code,
			"name":       county.Name,
			"created_at": county.CreatedAt,
		},
		// Cast the raw database bytes directly to json.RawMessage
		Geometry: json.RawMessage(county.Geometry),
	}

	// 4. Populate the Cache in the background (asynchronously)
	if featureBytes, err := json.Marshal(feature); err == nil {
		go func(ctx context.Context, key string, b []byte) {
			_ = s.cache.Set(ctx, key, b, 24*time.Hour)
		}(ctx, cacheKey, featureBytes)
	}

	return feature, nil
}

// ListCountiesAsFeatureCollection fetches all counties and packages them for Leaflet.js.
func (s *CountyService) ListCountiesAsFeatureCollection(ctx context.Context) (*geojson.FeatureCollection, error) {
	counties, err := s.repo.ListCounties(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list counties: %w", err)
	}

	features := make([]geojson.Feature, 0, len(counties))
	for _, county := range counties {
		feature := geojson.Feature{
			Type: "Feature",
			ID:   fmt.Sprintf("%d", county.ID),
			Properties: map[string]interface{}{
				"code": county.Code,
				"name": county.Name,
			},
			Geometry: json.RawMessage(county.Geometry),
		}
		features = append(features, feature)
	}

	collection := geojson.NewFeatureCollection(features)
	return collection, nil
}

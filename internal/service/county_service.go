package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
	"github.com/paula-dot/kenya-admin-boundaries-api/pkg/geojson"
)

// CountyRepository defines the required database operations.
// The PostgreSQL implementation will satisfy this interface.
type CountyRepository interface {
	GetCountyByCode(ctx context.Context, code string) (*domain.County, error)
	ListCounties(ctx context.Context) ([]*domain.County, error)
	GetCountyMetadata(ctx context.Context, code string) (postgres.GetCountyMetadataRow, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

// SpatialResult is the named return type used by SpatialIntersect.
// Using a named type allows runtime type assertions in handlers to succeed.
type SpatialResult struct {
	Ward         *domain.Ward
	Constituency *domain.Constituency
	County       *domain.County
}

// CountyService orchestrates business logic and data formatting.
type CountyService struct {
	repo        CountyRepository
	cache       CacheRepository
	spatialRepo repository.SpatialRepo
}

// NewCountyService is the constructor function.
func NewCountyService(repo CountyRepository, cache CacheRepository, spatial repository.SpatialRepo) *CountyService {
	return &CountyService{
		repo:        repo,
		cache:       cache,
		spatialRepo: spatial,
	}
}

// GetCountyAsFeature fetches a county by its official code and formats it as a standard GeoJSON Feature.
func (s *CountyService) GetCountyAsFeature(ctx context.Context, code string) (*geojson.Feature, error) {
	cacheKey := fmt.Sprintf("county:feature:%s", code)

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
	county, err := s.repo.GetCountyByCode(ctx, code)
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

// GetCountyBySlug is an adapter that allows the router to perform a slug/code lookup.
// We prefer using the spatial repository (which queries the `slug` column) when
// available, falling back to the legacy code-based lookup if necessary.
func (s *CountyService) GetCountyBySlug(ctx context.Context, slug string) (*domain.County, error) {
	// Try spatialRepo which performs slug-based lookup and returns a GeoJSON feature
	if s.spatialRepo != nil {
		feat, err := s.spatialRepo.GetCountyBySlug(ctx, slug)
		if err == nil && feat != nil {
			// Map geojson.Feature -> domain.County
			var id int32 = 0
			var codeStr string
			var nameStr string
			if feat.Properties != nil {
				if v, ok := feat.Properties["id"]; ok {
					switch t := v.(type) {
					case int:
						id = int32(t)
					case int32:
						id = t
					case float64:
						id = int32(t)
					}
				}
				if v, ok := feat.Properties["code"]; ok {
					if s, ok := v.(string); ok {
						codeStr = s
					}
				}
				if v, ok := feat.Properties["name"]; ok {
					if s, ok := v.(string); ok {
						nameStr = s
					}
				}
			}

			geometryBytes := []byte(nil)
			if len(feat.Geometry) > 0 {
				geometryBytes = []byte(feat.Geometry)
			}

			return &domain.County{
				ID:       id,
				Code:     codeStr,
				Name:     nameStr,
				Geometry: geometryBytes,
			}, nil
		}
	}

	// Fallback: attempt to treat the slug as an official code
	return s.repo.GetCountyByCode(ctx, slug)
}

// GetCountyMetadata provides lightweight code + name for a county.
func (s *CountyService) GetCountyMetadata(ctx context.Context, code string) (postgres.GetCountyMetadataRow, error) {
	return s.repo.GetCountyMetadata(ctx, code)
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

// SpatialIntersect implements the runtime interface used by the router.
// It returns up to Ward, Constituency and County wrapped in a named SpatialResult
// so handlers can perform runtime assertions against the service.
func (s *CountyService) SpatialIntersect(ctx context.Context, lat, lng float64) (SpatialResult, error) {
	var out SpatialResult
	if s.spatialRepo == nil {
		return out, fmt.Errorf("spatial repository not configured")
	}
	res, err := s.spatialRepo.GetLocationByPoint(ctx, lng, lat)
	if err != nil {
		return out, err
	}
	// Map plain names to domain objects. The repository returns names only; if
	// you have richer objects in DB consider extending the repo to return them.
	if res.Ward != "" {
		out.Ward = &domain.Ward{Name: res.Ward}
	}
	if res.Constituency != "" {
		out.Constituency = &domain.Constituency{Name: res.Constituency}
	}
	if res.County != "" {
		out.County = &domain.County{Name: res.County}
	}
	return out, nil
}

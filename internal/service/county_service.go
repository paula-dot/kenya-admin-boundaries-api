package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
	"github.com/paula-dot/kenya-admin-boundaries-api/pkg/geojson"
)

// CountyRepository defines the required database operations.
// The PostgreSQL implementation will satisfy this interface.
type CountyRepository interface {
	GetCountyByID(ctx context.Context, id uuid.UUID) (*domain.County, error)
	ListCounties(ctx context.Context) ([]*domain.County, error)
}

// CountyService orchestrates business logic and data formatting.
type CountyService struct {
	repo CountyRepository
}

// NewCountyService is the constructor function.
func NewCountyService(repo CountyRepository) *CountyService {
	return &CountyService{
		repo: repo,
	}
}

// GetCountyAsFeature fetches a county and formats it as a standard GeoJSON Feature.
func (s *CountyService) GetCountyAsFeature(ctx context.Context, id uuid.UUID) (*geojson.Feature, error) {
	county, err := s.repo.GetCountyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get county: %w", err)
	}

	feature := &geojson.Feature{
		Type: "Feature",
		ID:   county.ID.String(),
		Properties: map[string]interface{}{
			"code":       county.Code,
			"name":       county.Name,
			"created_at": county.CreatedAt,
		},
		// Cast the raw database bytes directly to json.RawMessage
		Geometry: json.RawMessage(county.Geometry),
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
			ID:   county.ID.String(),
			Properties: map[string]interface{}{
				"code": county.Code,
				"name": county.Name,
			},
			Geometry: json.RawMessage{county.Geometry},
		}
		features = append(features, feature)
	}

	collection := geojson.NewFeatureCollection(features)
	return collection, nil
}

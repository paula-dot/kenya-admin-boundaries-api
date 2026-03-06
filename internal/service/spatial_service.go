package service

import (
	"context"
	"errors"

	db "github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
)

// SpatialService acts as the bridge between the Gin handler and the Postgres engine
type SpatialService struct {
	queries *db.Queries
}

// NewSpatialService initializes the service with database access
func NewSpatialService(q *db.Queries) *SpatialService {
	return &SpatialService{queries: q}
}

// IntersectResponse represents the clean JSON payload returned to the client
type IntersectResponse struct {
	CountyCode       string `json:"county_code"`
	CountyName       string `json:"county_name"`
	ConstituencyCode string `json:"constituency_code"`
	ConstituencyName string `json:"constituency_name"`
}

// GetIntersection transforms the raw coordinates into a location payload
func (s *SpatialService) GetIntersection(ctx context.Context, lat, lon float64) (*IntersectResponse, error) {
	// Security & Safety: Basic bounds validation to prevent database strain
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return nil, errors.New("invalid coordinates: latitude must be between -90 and 90, and longitude between -180 and 180")
	}

	// PostGIS uses (Longitude, Latitude) for points (ST_MakePoint)
	params := db.GetIntersectingBoundaryParams{
		Longitude: lon,
		Latitude:  lat,
	}

	row, err := s.queries.GetIntersectingBoundary(ctx, params)
	if err != nil {
		// In production, log the exact error, but return a clean error to the user
		return nil, err
	}

	return &IntersectResponse{
		CountyCode:       row.CountyCode,
		CountyName:       row.CountyName,
		ConstituencyCode: row.ConstituencyCode,
		ConstituencyName: row.ConstituencyName,
	}, nil
}

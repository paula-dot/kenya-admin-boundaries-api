package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	db "github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
)

// SpatialService acts as the bridge between the Gin handler and the Postgres engine
type SpatialService struct {
	queries *db.Queries
	redis   *redis.Client
}

// NewSpatialService initializes the service with database access
func NewSpatialService(q *db.Queries, rdb *redis.Client) *SpatialService {
	return &SpatialService{
		queries: q,
		redis:   rdb,
	}
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

// GetLocationByCoordinates checks the cache first, then falls back to PostGIS
// and returns a generic FeatureCollection-like structure (interface{}) so handlers
// can marshal as GeoJSON if needed.
func (s *SpatialService) GetLocationByCoordinates(ctx context.Context, lon, lat float64) (interface{}, error) {
	// 1. Generate Cache Key (Round to 4 decimal places for ~11m precision)
	cacheKey := fmt.Sprintf("spatial:lookup:%.4f:%.4f", lon, lat)

	// 2. Try to fetch from Redis Cache
	cachedResult, err := s.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// CACHE MISS: Proceed to query PostGIS
		params := db.GetIntersectingBoundaryParams{Longitude: lon, Latitude: lat}
		row, err := s.queries.GetIntersectingBoundary(ctx, params)
		if err != nil {
			return nil, err
		}

		// Build a minimal FeatureCollection-like structure from the DB row
		dbResult := map[string]interface{}{
			"type": "FeatureCollection",
			"features": []interface{}{
				map[string]interface{}{
					"county_code":       row.CountyCode,
					"county_name":       row.CountyName,
					"constituency_code": row.ConstituencyCode,
					"constituency_name": row.ConstituencyName,
				},
			},
		}

		// Serialize the PostGIS result to JSON
		resultJSON, marshalErr := json.Marshal(dbResult)
		if marshalErr != nil {
			return nil, fmt.Errorf("failed to marshal db result: %w", marshalErr)
		}

		// 3. Save to Redis with a 24-hour TTL (Time To Live)
		err = s.redis.Set(ctx, cacheKey, resultJSON, 24*time.Hour).Err()
		if err != nil {
			// Log the error, but don't fail the request (cache writing shouldn't break the API)
			fmt.Printf("Warning: Failed to set cache for %s: %v\n", cacheKey, err)
		}

		return dbResult, nil

	} else if err != nil {
		// Handle actual Redis connection errors
		return nil, fmt.Errorf("redis error: %w", err)
	}

	// CACHE HIT: Unmarshal and return the cached data instantly
	var hitResult interface{}
	if err := json.Unmarshal([]byte(cachedResult), &hitResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached result: %w", err)
	}

	return hitResult, nil
}

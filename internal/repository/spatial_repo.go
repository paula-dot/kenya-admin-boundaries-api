package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paula-dot/kenya-admin-boundaries-api/pkg/geojson"
)

type SpatialRepository struct {
	DB *pgxpool.Pool
}

// IntersectionResult holds the plain JSON metadata of the matched location
type IntersectionResult struct {
	Ward         string `json:"ward"`
	Constituency string `json:"constituency"`
	County       string `json:"county"`
}

// GetCountyBySlug fetches a single county and formats it as a GeoJSON Feature
func (r *SpatialRepository) GetCountyBySlug(ctx context.Context, slug string) (*geojson.Feature, error) {
	query := `
		SELECT 
			'Feature' AS type,
			ST_AsGeoJSON(geom)::jsonb AS geometry,
			json_build_object(
				'id', id,
				'name', name,
				'slug', slug,
				'code', code
			) AS properties
		FROM counties
		WHERE slug = $1;
	`

	var feature geojson.Feature

	// Execute query and scan directly into the struct
	err := r.DB.QueryRow(ctx, query, slug).Scan(
		&feature.Type,
		&feature.Geometry,
		&feature.Properties,
	)

	if err != nil {
		return nil, err
	}
	return &feature, nil
}

// GetLocationByPoint checks which boundaries intersect with a given Lat/Lng
func (r *SpatialRepository) GetLocaationByPoint(ctx context.Context, lng, lat float64) (*IntersectionResult, error) {
	// ST_MakePoint takes (Longitude, Latitude). ST_SetSRID sets it to WGS84 standard (4326).
	query := `
		SELECT 
			  	w.name AS ward,
				c.name AS constituency,
				co.name AS county
		FROM wards w
		JOIN constituencies c ON w.constituency_id = c.id
		JOIN counties co ON c.county_id = co.id
		WHERE ST_Intersects(w.geom, ST_SetSRID(ST_MakePoint($1, $2), 4326));
	`
	var result IntersectionResult
	err := r.DB.QueryRow(ctx, query, lng, lat).Scan(&result.Ward, &result.Constituency, &result.County)

	if err != nil {
		return nil, fmt.Errorf("no boundaries found for these coordinates or database error: %w", err)
	}
	return &result, nil
}

package geojson

import "encoding/json"

// FeatureCollection represents a standard GeoJSON FeatureCollection.
// Leaflet.js will consume this directly to plot multiple layers on your map.
type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// Feature represents a standard GeoJSON Feature.
type Feature struct {
	Type       string                 `json:"type"`
	ID         string                 `json:"id,omitempty"`
	Properties map[string]interface{} `json:"properties"`

	// json.RawMessage tells Go to treat this as pre-formatted JSON.
	// This perfectly catches the `jsonb` output from our ST_AsGeoJSON SQL query.
	Geometry json.RawMessage `json:"geometry"`
}

// NewFeatureCollection is a helper function to quickly initialize a valid collection.
func NewFeatureCollection(features []Feature) *FeatureCollection {
	if features == nil {
		features = []Feature{} // Ensures it marshals to [] instead of null
	}

	return &FeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}
}

// --- Domain metadata structs for Feature.Properties ---
// CountyProperties holds the metadata for a county.
type CountyProperties struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Code string `json:"code,omitempty"` // E.g., "047" for Nairobi
}

// ConstituencyProperties holds metadata for a constituency.
type ConstituencyProperties struct {
	ID       int    `json:"id"`
	CountyID int    `json:"county_id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
}

// WardProperties holds metadata for a ward.
type WardProperties struct {
	ID             int    `json:"id"`
	ConstituencyID int    `json:"constituency_id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
}

// PointRequest represents the JSON body expected for spatial intersection queries.
type PointRequest struct {
	Longitude float64 `json:"lng" binding:"required,min=-180,max=180"`
	Latitude  float64 `json:"lat" binding:"required,min=-90,max=90"`
}

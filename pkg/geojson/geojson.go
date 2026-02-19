package geojson

import "encoding/json"

// FeatureCollection represents a standard GeoJSON FeatureCollection.
// Leaflet.js will consume this directly to plot multiple layers on your map.
type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
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

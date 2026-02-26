package dto

// PointRequest represents a simple lat/lng payload used by the spatial intersect endpoint.
// JSON keys are `lat` and `lng` to match the router's expected payload.
type PointRequest struct {
	Latitude  float64 `json:"lat" binding:"required"`
	Longitude float64 `json:"lng" binding:"required"`
}

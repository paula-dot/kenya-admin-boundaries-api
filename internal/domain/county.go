package domain

import "time"

type County struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Code      string    `json:"code,omitempty"`
	GeoJSON   string    `json:"geojson,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

package domain

import "time"

type Constituency struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CountyID  int       `json:"county_id"`
	Slug      string    `json:"slug"`
	GeoJSON   string    `json:"geojson,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

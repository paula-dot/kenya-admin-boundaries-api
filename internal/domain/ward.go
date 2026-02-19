package domain

import "time"

type Ward struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	ConstituencyID int       `json:"constituency_id"`
	Slug           string    `json:"slug"`
	GeoJSON        string    `json:"geojson,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

package domain

import (
	"time"
)

// County represents the top-level administrative boundary
type County struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Code      string    `json:"code,omitempty"`
	GeoJSON   string    `json:"geojson,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// SubCounty represents administrative areas under a county (No spatial data
type SubCounty struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CountyID  int       `json:"county_id"`
	Slug      string    `json:"slug"`
	GeoJSON   string    `json:"geojson,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Constituency represents electoral boundaries within a county
type Constituency struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CountyID  int    `json:"county_id"`
	Slug      string `json:"slug"`
	GeoJSON   string `json:"geojson,omitempty"`
	CreatedAt time.Time
}

type Ward struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ConstituencyID int    `json:"constituency_id"`
	Slug           string `json:"slug"`
	GeoJSON        string `json:"geojson,omitempty"`
	CreatedAt      time.Time
}

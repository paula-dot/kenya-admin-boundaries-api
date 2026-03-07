package domain

import (
	"time"
)

// County represents the core administrative boundary entity.
type County struct {
	ID        int32
	Code      string
	Name      string
	Geometry  []byte // This will hold the raw JSON bytes returned from the DB
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SubCounty represents an administrative unit within a County.
type SubCounty struct {
	CountyCode    string `json:"county_code"`
	CountyName    string `json:"county_name"`
	SubCountyCode string `json:"sub_county_code"`
	SubCountyName string `json:"sub_county_name"`
}

package domain

import (
	"time"
)

type GeoJSON string

type Base struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

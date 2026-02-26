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

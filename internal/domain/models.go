package domain

import (
	"time"

	"github.com/google/uuid"
)

// County represents the core administrative boundary entity.
type County struct {
	ID        uuid.UUID
	Code      string
	Name      string
	Geometry  []byte // This will hold the raw JSON bytes returned from the DB
	CreatedAt time.Time
	UpdatedAt time.Time
}

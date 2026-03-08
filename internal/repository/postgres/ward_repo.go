package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// WardRepository encapsulates database access for Wards
type WardRepository struct {
	db *pgxpool.Pool
	q  *Queries
}

// NewWardRepository creates a new repository using the connection pool
func NewWardRepository(db *pgxpool.Pool) *WardRepository {
	return &WardRepository{
		db: db,
		q:  New(db),
	}
}

// ListWards fetches a paginated list of wards.
func (r *WardRepository) ListWards(ctx context.Context, limit, offset int32) ([]ListWardsRow, error) {
	return r.q.ListWards(ctx, ListWardsParams{
		Limit:  limit,
		Offset: offset,
	})
}

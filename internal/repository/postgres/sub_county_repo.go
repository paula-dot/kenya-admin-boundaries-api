package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
)

// SubCountyRepo implements repository methods for sub-counties.
type SubCountyRepo struct {
	db Querier
}

// NewSubCountyRepository creates a new PostgreSQL-backed sub-county repository.
func NewSubCountyRepository(pool *pgxpool.Pool) *SubCountyRepo {
	return &SubCountyRepo{db: New(pool)}
}

// GetAllSubCounties retrieves all sub-counties across the country.
func (r *SubCountyRepo) GetAllSubCounties(ctx context.Context) ([]*domain.SubCounty, error) {
	rows, err := r.db.GetAllSubCounties(ctx)
	if err != nil {
		return nil, err
	}

	var results []*domain.SubCounty
	for _, row := range rows {
		results = append(results, &domain.SubCounty{
			CountyCode:    row.CountyCode,
			CountyName:    row.CountyName,
			SubCountyCode: row.SubCountyCode,
			SubCountyName: row.SubCountyName,
		})
	}
	return results, nil
}

// GetSubCountiesByCounty retrieves all sub-counties for a specific county code.
func (r *SubCountyRepo) GetSubCountiesByCounty(ctx context.Context, countyCode string) ([]*domain.SubCounty, error) {
	rows, err := r.db.GetSubCountiesByCounty(ctx, countyCode)
	if err != nil {
		return nil, err
	}

	var results []*domain.SubCounty
	for _, row := range rows {
		results = append(results, &domain.SubCounty{
			CountyCode:    countyCode, // The query doesn't explicitly return this, we have it from input
			CountyName:    row.CountyName,
			SubCountyCode: row.SubCountyCode,
			SubCountyName: row.SubCountyName,
		})
	}
	return results, nil
}

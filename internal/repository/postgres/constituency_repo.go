package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
)

// ConstituencyRepo implements repository methods for constituencies.
type ConstituencyRepo struct {
	db Querier
}

// NewConstituencyRepository creates a new PostgreSQL-backed constituency repository.
func NewConstituencyRepository(pool *pgxpool.Pool) *ConstituencyRepo {
	return &ConstituencyRepo{db: New(pool)}
}

// GetConstituencyByCode executes the SQL query and maps the result to the domain model.
func (r *ConstituencyRepo) GetConstituencyByCode(ctx context.Context, code string) (*domain.Constituency, error) {
	row, err := r.db.GetConstituencyByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("repository - GetConstituencyByCode failed: %w", err)
	}

	cCode := ""
	if row.ID != "" {
		cCode = row.ID
	}

	created := time.Time{}

	var idVal int = 0
	if parsed, err := strconv.Atoi(cCode); err == nil {
		idVal = parsed
	}

	var countyID int = 0
	if row.CountyCode != "" {
		if parsed, err := strconv.Atoi(row.CountyCode); err == nil {
			countyID = parsed
		}
	}

	cons := &domain.Constituency{
		ID:        idVal,
		Name:      row.Name,
		CountyID:  countyID,
		Slug:      cCode,
		GeoJSON:   string(row.Geojson),
		CreatedAt: created,
	}

	return cons, nil
}

// ListConstituencies executes the list query and maps the slice of results to domain models.
func (r *ConstituencyRepo) ListConstituencies(ctx context.Context) ([]*domain.Constituency, error) {
	rows, err := r.db.ListConstituencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository - ListConstituencies failed: %w", err)
	}

	var out []*domain.Constituency
	for _, row := range rows {
		slug := ""
		if row.ID != "" {
			slug = row.ID
		}

		var idVal int = 0
		if parsed, err := strconv.Atoi(slug); err == nil {
			idVal = parsed
		}

		var countyID int = 0
		if row.CountyCode != "" {
			if parsed, err := strconv.Atoi(row.CountyCode); err == nil {
				countyID = parsed
			}
		}

		c := &domain.Constituency{
			ID:       idVal,
			Name:     row.Name,
			CountyID: countyID,
			Slug:     slug,
			GeoJSON:  string(row.Geojson),
		}
		out = append(out, c)
	}

	return out, nil
}

// ListConstituenciesByCounty returns the raw sqlc rows for a county_code.
// The service layer will transform these into the shape expected by handlers.
func (r *ConstituencyRepo) ListConstituenciesByCounty(ctx context.Context, countyCode string) ([]ListConstituenciesByCountyRow, error) {
	return r.db.ListConstituenciesByCounty(ctx, countyCode)
}

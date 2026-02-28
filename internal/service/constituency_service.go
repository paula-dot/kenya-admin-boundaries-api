package service

import (
	"context"
	"strconv"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
)

// ConstituencyListItem is the shape expected by handlers/router for constituency lists.
type ConstituencyListItem struct {
	ID       int32
	Slug     string
	Name     string
	Geometry []byte
}

// constituencyRepo defines the subset of repository methods used by the service.
type constituencyRepo interface {
	ListConstituenciesByCounty(ctx context.Context, countyCode string) ([]postgres.ListConstituenciesByCountyRow, error)
}

// ConstituencyService provides helper methods around constituencies.
type ConstituencyService struct {
	repo constituencyRepo
}

// NewConstituencyService constructs a new service.
func NewConstituencyService(repo constituencyRepo) *ConstituencyService {
	return &ConstituencyService{repo: repo}
}

// ListConstituenciesByCountySlug returns constituencies for a given county code/slug.
func (s *ConstituencyService) ListConstituenciesByCountySlug(ctx context.Context, slug string) ([]struct {
	ID       int32
	Slug     string
	Name     string
	Geometry []byte
}, error) {
	rows, err := s.repo.ListConstituenciesByCounty(ctx, slug)
	if err != nil {
		return nil, err
	}

	out := make([]struct {
		ID       int32
		Slug     string
		Name     string
		Geometry []byte
	}, 0, len(rows))

	for _, r := range rows {
		var id int32 = 0
		if r.ID != "" {
			if parsed, err := strconv.Atoi(r.ID); err == nil {
				id = int32(parsed)
			}
		}
		item := struct {
			ID       int32
			Slug     string
			Name     string
			Geometry []byte
		}{
			ID:       id,
			Slug:     r.ID,
			Name:     r.Name,
			Geometry: r.Geojson,
		}
		out = append(out, item)
	}

	return out, nil
}

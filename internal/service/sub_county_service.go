package service

import (
	"context"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
)

// subCountyRepo defines the methods we need from the database package.
type subCountyRepo interface {
	GetAllSubCounties(ctx context.Context) ([]*domain.SubCounty, error)
	GetSubCountiesByCounty(ctx context.Context, countyCode string) ([]*domain.SubCounty, error)
}

// SubCountyService handles business logic related to Sub-Counties.
type SubCountyService struct {
	repo subCountyRepo
}

// NewSubCountyService initializes a new sub-county service.
func NewSubCountyService(repo subCountyRepo) *SubCountyService {
	return &SubCountyService{repo: repo}
}

// GetAll returns a complete list of all sub-counties.
func (s *SubCountyService) GetAll(ctx context.Context) ([]*domain.SubCounty, error) {
	return s.repo.GetAllSubCounties(ctx)
}

// GetByCountyCode returns sub-counties filtered by a given county code.
func (s *SubCountyService) GetByCountyCode(ctx context.Context, code string) ([]*domain.SubCounty, error) {
	return s.repo.GetSubCountiesByCounty(ctx, code)
}

package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
)

// fakeRepo implements the constituencyRepo interface for tests.
type fakeRepo struct {
	rows []postgres.ListConstituenciesByCountyRow
	err  error
}

func (f *fakeRepo) ListConstituenciesByCounty(ctx context.Context, countyCode string) ([]postgres.ListConstituenciesByCountyRow, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.rows, nil
}

func (f *fakeRepo) ListConstituenciesMetadataByCounty(ctx context.Context, countyCode string) ([]postgres.ListConstituenciesMetadataByCountyRow, error) {
	return nil, f.err
}

func TestListConstituenciesByCountySlug(t *testing.T) {
	tests := []struct {
		name    string
		rows    []postgres.ListConstituenciesByCountyRow
		repoErr error
		want    []struct {
			ID       int32
			Slug     string
			Name     string
			Geometry []byte
		}
		wantErr bool
	}{
		{
			name: "happy-path parses numeric ids",
			rows: []postgres.ListConstituenciesByCountyRow{
				{ID: "1", Name: "One", CountyCode: "10", Geojson: []byte(`{"type":"Point"}`)},
				{ID: "2", Name: "Two", CountyCode: "10", Geojson: []byte(`{"type":"Point"}`)},
			},
			want: []struct {
				ID       int32
				Slug     string
				Name     string
				Geometry []byte
			}{
				{ID: 1, Slug: "1", Name: "One", Geometry: []byte(`{"type":"Point"}`)},
				{ID: 2, Slug: "2", Name: "Two", Geometry: []byte(`{"type":"Point"}`)},
			},
			wantErr: false,
		},
		{
			name: "non-numeric id defaults to zero",
			rows: []postgres.ListConstituenciesByCountyRow{
				{ID: "A-100", Name: "Alpha", CountyCode: "10", Geojson: []byte(`null`)},
			},
			want: []struct {
				ID       int32
				Slug     string
				Name     string
				Geometry []byte
			}{
				{ID: 0, Slug: "A-100", Name: "Alpha", Geometry: []byte(`null`)},
			},
			wantErr: false,
		},
		{
			name:    "repo returns error",
			rows:    nil,
			repoErr: errors.New("db down"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{rows: tc.rows, err: tc.repoErr}
			svc := NewConstituencyService(repo)
			got, err := svc.ListConstituenciesByCountySlug(context.Background(), "10")
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error state: got err=%v wantErr=%v", err, tc.wantErr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("unexpected result:\n got: %#v\nwant: %#v", got, tc.want)
			}
		})
	}
}

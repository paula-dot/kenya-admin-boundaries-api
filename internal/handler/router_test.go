package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
	"github.com/paula-dot/kenya-admin-boundaries-api/pkg/geojson"
)

// mockService implements the minimal methods that SetupRouter performs runtime
// type assertions against: GetCountyBySlug and SpatialIntersect.
type mockService struct {
	// store values for assertions if needed
}

func (m *mockService) ListCountiesAsFeatureCollection(ctx context.Context) (*geojson.FeatureCollection, error) {
	return &geojson.FeatureCollection{Type: "FeatureCollection", Features: []geojson.Feature{}}, nil
}

func (m *mockService) GetCountyBySlug(ctx context.Context, slug string) (*domain.County, error) {
	// return a simple county with geometry bytes
	return &domain.County{
		ID:       1,
		Name:     "TestCounty",
		Code:     "TC",
		Geometry: []byte(`{"type":"Point","coordinates":[0,0]}`),
	}, nil
}

func (m *mockService) SpatialIntersect(ctx context.Context, lat, lng float64) (struct {
	Ward         *domain.Ward
	Constituency *domain.Constituency
	County       *domain.County
}, error) {
	return struct {
		Ward         *domain.Ward
		Constituency *domain.Constituency
		County       *domain.County
	}{
		Ward:         nil,
		Constituency: nil,
		County:       &domain.County{ID: 1, Name: "TestCounty", Geometry: []byte(`{"type":"Point","coordinates":[0,0]}`)},
	}, nil
}

func TestGetCountyBySlugRoute(t *testing.T) {
	m := &mockService{}
	r := SetupRouter(m)

	req := httptest.NewRequest("GET", "/api/v1/counties/test-slug", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}

	// Expect GeoJSON FeatureCollection in response
	if ct := w.Header().Get("Content-Type"); ct != "application/geo+json" {
		t.Fatalf("expected content-type application/geo+json got %s", ct)
	}
}

func TestSpatialIntersectRoute(t *testing.T) {
	m := &mockService{}
	r := SetupRouter(m)

	payload := map[string]float64{"lat": -1.0, "lng": 36.0}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/spatial/intersect", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/geo+json" {
		t.Fatalf("expected content-type application/geo+json got %s", ct)
	}
}

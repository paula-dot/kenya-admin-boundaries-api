package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/domain"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"
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

func (m *mockService) SpatialIntersect(ctx context.Context, lat, lng float64) (service.SpatialResult, error) {
	return service.SpatialResult{
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

type mockConstSvc struct{}

func (m *mockConstSvc) ListConstituenciesByCountySlug(ctx context.Context, slug string) ([]struct {
	ID       int32
	Slug     string
	Name     string
	Geometry []byte
}, error) {
	return []struct {
		ID       int32
		Slug     string
		Name     string
		Geometry []byte
	}{
		{ID: 123, Slug: "KE001-01", Name: "Test Constituency", Geometry: []byte(`{"type":"Point","coordinates":[36,-1]}`)},
	}, nil
}

func TestListConstituenciesByCountyRoute(t *testing.T) {
	m := &mockConstSvc{}
	r := SetupRouter(m)

	req := httptest.NewRequest("GET", "/api/v1/counties/KE001/constituencies", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d: %s", w.Code, w.Body.String())
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/geo+json" {
		t.Fatalf("expected content-type application/geo+json got %s", ct)
	}

	// quick body sanity check: should contain "FeatureCollection" and the name
	b := w.Body.String()
	if b == "" {
		t.Fatalf("empty body returned")
	}
	if !strings.Contains(b, "FeatureCollection") || !strings.Contains(b, "Test Constituency") {
		t.Fatalf("unexpected body: %s", b)
	}
}

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/config"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/handler"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

// noopCache implements the minimal CacheRepository used by CountyService in tests.
type noopCache struct{}

func (n *noopCache) Get(ctx context.Context, key string) ([]byte, error) { return nil, nil }
func (n *noopCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

// LocationResponse is a minimal struct used to decode the handler response (if any).
type LocationResponse struct {
	CountyCode       string `json:"county_code"`
	CountyName       string `json:"county_name"`
	ConstituencyCode string `json:"constituency_code,omitempty"`
	ConstituencyName string `json:"constituency_name,omitempty"`
}

func TestLocationEndpoint(t *testing.T) {
	// Only run when explicitly requested (avoid slow integration by default)
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("skipping integration test; set INTEGRATION=true to enable")
	}

	// Load config (allows using .env in repo root)
	cfg, _ := config.LoadConfig(".")
	dsn := cfg.DBUrl
	if dsn == "" {
		// Fallback to common local development mapping (override maps to 127.0.0.1:5433)
		dsn = "postgres://user:Kapcherop250.@127.0.0.1:5433/spatial_db?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Skipf("skipping integration test; unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Quick ping to ensure DB is responsive
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("skipping integration test; database not responding to ping: %v", err)
	}

	// Build real repository/service (uses noop cache so redis is not required)
	pgRepo := postgres.NewCountyRepository(pool)
	svc := service.NewCountyService(pgRepo, &noopCache{})

	// Build router from the project's handler.SetupRouter
	r := handler.SetupRouter(svc)

	// Ensure the reverse-geocoding route exists at /api/v1/location by registering it here
	api := r.Group("/api/v1")
	api.GET("/location", handler.GetLocationByPoint(pool))

	// Start an in-process HTTP server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Example coordinates (replace with coordinates that fall within a known area if needed)
	url := ts.URL + "/api/v1/location?lat=-0.514&lng=35.269"

	res, err := http.Get(url)
	if err != nil {
		t.Fatalf("http get failed: %v", err)
	}
	defer res.Body.Close()

	// Fail only on server errors (500+). 404 is acceptable for a template run.
	if res.StatusCode >= 500 {
		t.Fatalf("server error: status %d", res.StatusCode)
	}

	// Optionally decode JSON response when status 200
	if res.StatusCode == 200 {
		var lr LocationResponse
		if err := json.NewDecoder(res.Body).Decode(&lr); err != nil {
			t.Fatalf("failed to decode response JSON: %v", err)
		}
		// Basic assertion: county name or code should be present
		if lr.CountyName == "" && lr.CountyCode == "" {
			t.Fatalf("empty location response")
		}
	}
}

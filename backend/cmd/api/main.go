package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/config"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/handler"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/middleware"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
	redisRepo "github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/redis"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// noopCache is used when Redis is not configured or fails to initialize in dev.
// It implements the CacheRepository interface used by service.NewCountyService.
type noopCache struct{}

func (n *noopCache) Get(ctx context.Context, key string) ([]byte, error) { return nil, nil }
func (n *noopCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

func main() {
	// 1. Initialize Context
	ctx := context.Background()

	// 2. Load Configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Unable to load configuration: %v\n", err)
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DBUrl == "" {
		// keep a sensible fallback for environments where .env is not present
		cfg.DBUrl = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	// Print final DSN but redact the password when possible
	if u, err := url.Parse(cfg.DBUrl); err == nil {
		if u.User != nil {
			username := u.User.Username()
			// replace password with REDACTED for safe logging
			u.User = url.UserPassword(username, "REDACTED")
		}
		log.Printf("Using DATABASE_URL: %s\n", u.String())
	} else {
		log.Printf("Using DATABASE_URL (unparsed): %s\n", cfg.DBUrl)
	}

	// 3. Database Connection Pooling (pgxpool)
	dbPool, err := pgxpool.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("Unable to create database connection pool: %v\n", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	log.Println("Successfully connected to database!")

	// Create sqlc-generated Queries (used by SpatialService)
	sqlcQueries := postgres.New(dbPool)

	pgRepo := postgres.NewCountyRepository(dbPool)
	// create constituency repository and service so router can list constituencies
	consRepo := postgres.NewConstituencyRepository(dbPool)
	// create sub-county repository
	subCountyRepo := postgres.NewSubCountyRepository(dbPool)
	// create ward repository
	wardRepo := postgres.NewWardRepository(dbPool)

	// Initialize redis cache if configured; otherwise fall back to noopCache
	var cacheRepo service.CacheRepository
	// Also create a redis.Client for the SpatialService. Parse cfg.RedisURL if present.
	var rdb *redis.Client
	if cfg.RedisURL == "" {
		log.Println("REDIS_URL not set; using noop cache (development mode)")
		cacheRepo = &noopCache{}
		// create a local redis client pointing at localhost for best-effort caching in dev
		rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	} else {
		// Try to initialize both the repository wrapper and the raw redis client.
		r, err := redisRepo.NewCacheRepository(cfg.RedisURL)
		if err != nil {
			log.Printf("Unable to initialize redis cache repo: %v; falling back to noop cache\n", err)
			cacheRepo = &noopCache{}
		} else {
			cacheRepo = r
		}

		// Parse cfg.RedisURL to extract host:port and optional password for go-redis
		addr := cfg.RedisURL
		password := ""
		if u, err := url.Parse(cfg.RedisURL); err == nil && u.Host != "" {
			addr = u.Host
			if u.User != nil {
				if p, ok := u.User.Password(); ok {
					password = p
				}
			}
		}
		rdb = redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: 0})
	}

	// Attempt a lightweight ping to Redis; don't fail startup if Redis is unreachable
	if rdb != nil {
		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Printf("Warning: Redis ping failed: %v (continuing without cache client)\n", err)
		}
	}

	spatialRepo := &repository.SpatialRepository{DB: dbPool}

	countySvc := service.NewCountyService(pgRepo, cacheRepo, spatialRepo)
	consSvc := service.NewConstituencyService(consRepo, cacheRepo)
	spatialSvc := service.NewSpatialService(sqlcQueries, rdb)
	subCountySvc := service.NewSubCountyService(subCountyRepo, cacheRepo)
	wardSvc := service.NewWardService(wardRepo, cacheRepo)

	// Use the handler.AppServices type so SetupRouter can register explicit handlers
	svc := &handler.AppServices{
		County:       countySvc,
		Constituency: consSvc,
		Spatial:      spatialSvc,
		SubCounty:    subCountySvc,
		Ward:         wardSvc,
	}

	// wire svc into handlers/router and apply rate limiter middleware to /api/v1
	router := handler.SetupRouter(svc, middleware.RateLimiter(rdb, 60, time.Minute))

	// Health endpoint (ensure SetupRouter sets other routes; this is a safe fallback)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "API is up and running!",
			"env":    cfg.Environment,
		})
	})

	// Expose routes list for quick debugging
	router.GET("/__routes", func(c *gin.Context) {
		routes := router.Routes()
		type simpleRoute struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		}
		out := make([]simpleRoute, 0, len(routes))
		for _, r := range routes {
			out = append(out, simpleRoute{Method: r.Method, Path: r.Path})
		}
		c.JSON(http.StatusOK, out)
	})

	// Note: route registration happens inside SetupRouter. Avoid registering the
	// same routes again here to prevent duplicate route registration panics.

	// 5. Server Initialization
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// 6. Graceful Shutdown Implementation
	go func() {
		log.Printf("Starting API server on port %s...\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v\n", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received. Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	if rdb != nil {
		rdb.Close()
	}

	log.Println("Server shut down successfully")
}

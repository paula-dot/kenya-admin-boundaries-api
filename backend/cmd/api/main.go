package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/config"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/handler"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/postgres"
	redisRepo "github.com/paula-dot/kenya-admin-boundaries-api/internal/repository/redis"
	"github.com/paula-dot/kenya-admin-boundaries-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/paula-dot/kenya-admin-boundaries-api/internal/repository"
	"github.com/paula-dot/kenya-admin-boundaries-api/pkg/db"
)

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
		cfg.DBUrl = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
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

	// 4. Dependency Injection Setup
	pgRepo := postgres.NewCountyRepository(dbPool)

	cacheRepo, err := redisRepo.NewCacheRepository(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Unable to initialize redis cache: %v\n", err)
	}

	svc := service.NewCountyService(pgRepo, cacheRepo)

	// wire svc into handlers/router
	router := handler.SetupRouter(svc)

	// Health endpoint (ensure SetupRouter sets other routes; this is a safe fallback)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "API is up and running!",
			"env":    cfg.Environment,
		})
	})

	// Additional spatial repository and handler wiring (user requested append)
	// Inside main()
	// NOTE: these lines assume packages `pkg/db`, `internal/repository` and `internal/handler`
	// expose the referenced symbols. Adjust imports/paths if your project differs.
	dbPool2 := db.Connect("postgresql://postgres:password@localhost:5433/spatial_db")
	repo := &repository.SpatialRepository{DB: dbPool2}
	apiHandler := &handler.APIHandler{Repo: repo}

	// If SetupRouter created route groups `v1` and `spatial` these show how to wire them.
	// Otherwise you may need to adapt to your router structure.
	v1 := router.Group("/v1")
	spatial := v1.Group("/spatial")
	v1.GET("/counties/:slug", apiHandler.GetCountyBySlug)
	spatial.POST("/intersect", apiHandler.CheckIntersection)

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
	log.Println("Server shut down successfully")
}

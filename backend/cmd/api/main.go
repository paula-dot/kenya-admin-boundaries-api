package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	// 1. Initialize Context
	ctx := context.Background()

	// 2. Load Configuration
	// TODO: Replace with your internal/config loader once built
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback for local development
		dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 3. Database Connection Pooling (pgxpool)
	// pgxpool manages a pool of connections securely and efficiently.
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to create database connection pool: %v\n", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	log.Println("Successfully connected to database!")

	// 4. Dependency Injection Setup
	// TODO: Instantiate Redis cache
	// TODO: repo := postgres.NewCountyRepository(dbPool)
	// TODO: svc := service.NewCountyService(repo)
	// TODO: router := handler.SetupRouter(svc)

	// Temporary Gin Router setup for testing
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API is up and running!"})
	})

	// 5. Server Initialization
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// 6. Graceful Shutdown Implementation
	// Run the server in a goroutine so it doesn't block the shutdown handling
	go func() {
		log.Printf("Starting API server on port %s...\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v\n", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
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

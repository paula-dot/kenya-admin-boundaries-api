package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect initializes and returns a connection pool
func Connect(databaseURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Verify connection is alive
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}

	log.Println("Successfully connected to PostGIS pool!")
	return pool
}

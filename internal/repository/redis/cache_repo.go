package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrCacheMiss is a custom error indicating the key was not found in Redis.
// This allows the service layer to gracefully fall back to the database.
var ErrCacheMiss = errors.New("cache: key not found")

// CacheRepo implements a generic Redis caching layer
type CacheRepo struct {
	client *redis.Client
}

func NewCacheRepository(redisURL string) (*CacheRepo, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis: unable to parse URL: %w", err)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: unable to connect: %w", err)
	}

	return &CacheRepo{
		client: client,
	}, nil
}

// Set stores a key-value pair in Redis with an expiration time.
// We use []byte for the value so we can directly store the json.RawMessage (GeoJSON).
func (r *CacheRepo) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis: unable to set key %s: %w", key, err)
	}
	return nil
}

// Get retrieves a value from Redis by key.
func (r *CacheRepo) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("redis: unable to get key %s: %w", key, err)
	}
	return val, nil
}

// Delete removes a key from Redis.
func (r *CacheRepo) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis: unable to delete key %s: %w", key, err)
	}
	return nil
}

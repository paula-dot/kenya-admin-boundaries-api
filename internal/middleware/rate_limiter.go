package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter creates a fixed-window rate limiter using Redis
func RateLimiter(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if rdb == nil || limit <= 0 {
			c.Next()
			return
		}

		// 1. Identify the user by their IP address
		clientIP := c.ClientIP()

		// Create a unique Redis key for this IP based on the current window
		windowKey := time.Now().Truncate(window).Unix()
		key := fmt.Sprintf("ratelimit:%s:%d", clientIP, windowKey)

		// 2. Increment the request counter for this time window
		count, err := rdb.Incr(c.Request.Context(), key).Result()
		if err != nil {
			// If Redis is down, we "fail open" so legitimate traffic isn't blocked
			fmt.Printf("Rate limiter Redis error: %v\n", err)
			c.Next()
			return
		}

		// 3. Set a TTL on the key if it's the first request in the window
		if count == 1 {
			_ = rdb.Expire(c.Request.Context(), key, window).Err()
		}

		// 4. Check if the limit has been exceeded
		if count > int64(limit) {
			ttl, terr := rdb.TTL(c.Request.Context(), key).Result()
			retryAfter := int(ttl.Seconds())
			if terr != nil || retryAfter < 0 {
				retryAfter = int(window.Seconds())
			}
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests. Please slow down.",
				"retry_after": retryAfter,
			})
			return
		}

		// (Optional) Inject helpful headers for the client
		remaining := int64(limit) - count
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", int(window.Seconds())))

		c.Next()
	}
}

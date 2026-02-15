package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

// CacheConfig holds cache configuration
type CacheConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	DefaultTTL    time.Duration
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		DefaultTTL:    5 * time.Minute, // Default cache TTL of 5 minutes
	}
}

// InitializeRedis initializes the Redis client
func InitializeRedis(config *CacheConfig) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Test connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis successfully")
	return nil
}

// GetRedisClient returns the Redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// GenerateCacheKey generates a cache key from the request path and query params
func GenerateCacheKey(c *gin.Context) string {
	return c.Request.URL.Path + "?" + c.Request.URL.RawQuery
}

// CacheMiddleware returns a middleware that caches GET responses
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Generate cache key
		cacheKey := GenerateCacheKey(c)

		// Try to get from cache
		cachedData, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			// Cache hit - return cached response
			log.Printf("Cache HIT for key: %s", cacheKey)

			// Try to unmarshal as map first, then as array
			var responseMap map[string]interface{}
			var responseArray []interface{}

			if err := json.Unmarshal([]byte(cachedData), &responseMap); err == nil {
				c.JSON(http.StatusOK, responseMap)
			} else if err := json.Unmarshal([]byte(cachedData), &responseArray); err == nil {
				c.JSON(http.StatusOK, responseArray)
			} else {
				log.Printf("Error unmarshaling cached response: %v", err)
				c.Next()
				return
			}
			c.Abort()
			return
		}

		// Cache miss - continue to handler and cache the response
		log.Printf("Cache MISS for key: %s", cacheKey)

		// Create a response writer wrapper to capture the response
		writer := &ResponseWriter{
			ResponseWriter: c.Writer,
			body:           nil,
			statusCode:     http.StatusOK,
		}
		c.Writer = writer

		c.Next()

		// Cache the response if it's successful
		if writer.statusCode >= http.StatusOK && writer.statusCode < http.StatusMultipleChoices {
			bodyBytes, err := json.Marshal(writer.body)
			if err != nil {
				log.Printf("Error marshaling response for cache: %v", err)
				return
			}

			err = redisClient.Set(ctx, cacheKey, bodyBytes, ttl).Err()
			if err != nil {
				log.Printf("Error setting cache: %v", err)
				return
			}
			log.Printf("Cached response for key: %s with TTL: %v", cacheKey, ttl)
		}
	}
}

// ResponseWriter is a wrapper around gin.ResponseWriter to capture the response body
type ResponseWriter struct {
	gin.ResponseWriter
	body       interface{}
	statusCode int
}

// Write captures the response body
func (rw *ResponseWriter) Write(body []byte) (int, error) {
	// Try to unmarshal and store the body
	var data interface{}
	if err := json.Unmarshal(body, &data); err == nil {
		rw.body = data
	}
	return rw.ResponseWriter.Write(body)
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// InvalidateCache deletes a specific cache key
func InvalidateCache(key string) error {
	return redisClient.Del(ctx, key).Err()
}

// InvalidateCachePattern deletes all keys matching a pattern
func InvalidateCachePattern(pattern string) error {
	iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return redisClient.Del(ctx, keys...).Err()
	}
	return nil
}

// InvalidateCacheByPrefix invalidates cache for a specific endpoint prefix
func InvalidateCacheByPrefix(prefix string) error {
	return InvalidateCachePattern(prefix + "*")
}

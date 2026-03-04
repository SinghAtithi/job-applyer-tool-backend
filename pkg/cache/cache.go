package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/pkg/logger"

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
		DefaultTTL:    5 * time.Minute,
	}
}

// InitializeRedis initializes the Redis client
func InitializeRedis(config *CacheConfig) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		// Reset client to nil so CacheMiddleware knows Redis is unavailable
		redisClient = nil
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis successfully")
	return nil
}

// GetRedisClient returns the Redis client (may be nil if not initialized)
func GetRedisClient() *redis.Client {
	return redisClient
}

// GenerateCacheKey generates a cache key from the request path and query params
func GenerateCacheKey(c *gin.Context) string {
	return c.Request.URL.Path + "?" + c.Request.URL.RawQuery
}

// CacheMiddleware returns a middleware that caches GET responses.
// If Redis is not available, requests pass through without caching.
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Guard: skip caching if Redis is not available
		if redisClient == nil {
			c.Next()
			return
		}

		cacheKey := GenerateCacheKey(c)

		// Try to get from cache
		cachedData, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			logger.Debug("cache HIT for key: %s", cacheKey)

			var responseMap map[string]interface{}
			var responseArray []interface{}

			if err := json.Unmarshal([]byte(cachedData), &responseMap); err == nil {
				c.JSON(http.StatusOK, responseMap)
			} else if err := json.Unmarshal([]byte(cachedData), &responseArray); err == nil {
				c.JSON(http.StatusOK, responseArray)
			} else {
				logger.Error("failed to unmarshal cached response for key %s: %v", cacheKey, err)
				c.Next()
				return
			}
			c.Abort()
			return
		}

		// Cache miss
		logger.Debug("cache MISS for key: %s", cacheKey)

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
				logger.Error("failed to marshal response for cache: %v", err)
				return
			}

			if setErr := redisClient.Set(ctx, cacheKey, bodyBytes, ttl).Err(); setErr != nil {
				logger.Error("failed to set cache for key %s: %v", cacheKey, setErr)
				return
			}
			logger.Debug("cached response for key: %s (TTL: %v)", cacheKey, ttl)
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
	if redisClient == nil {
		return nil
	}
	return redisClient.Del(ctx, key).Err()
}

// InvalidateCachePattern deletes all keys matching a pattern
func InvalidateCachePattern(pattern string) error {
	if redisClient == nil {
		return nil
	}
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

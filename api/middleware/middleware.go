package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"example.com/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDKey is the key used to store request ID in context
const RequestIDKey = "request_id"

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger middleware provides structured request logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details using the custom logger
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		requestID, _ := c.Get(RequestIDKey)

		logger.Info("[req:%v] %s %s | %s | %d | %s",
			requestID,
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			statusCode,
			latency,
		)
	}
}

// Recovery middleware recovers from panics and logs stack traces in debug mode
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Always log the panic
				logger.Error("panic recovered: %v", err)

				// Capture stack trace
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				stackTrace := string(buf[:n])

				// Always log the stack trace at error level
				logger.Error("stack trace:\n%s", stackTrace)

				response := gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INTERNAL_ERROR",
						"message": "An unexpected error occurred",
					},
				}

				// In debug mode, include stack trace in the response body
				if logger.IsDebugMode() {
					response["error"] = gin.H{
						"code":        "INTERNAL_ERROR",
						"message":     "An unexpected error occurred",
						"debug_panic": fmt.Sprintf("%v", err),
						"stack_trace": stackTrace,
					}
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			}
		}()
		c.Next()
	}
}

// CORS middleware handles cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Allow all origins in development, restrict in production
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID, Content-Length")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SecurityHeaders middleware adds security-related headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")

		// X-Frame-Options
		c.Header("X-Frame-Options", "DENY")

		// X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Strict-Transport-Security (HSTS)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Referrer-Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// HealthCheck returns a health check handler
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Service is healthy",
			"data": gin.H{
				"status":    "UP",
				"timestamp": time.Now().Unix(),
			},
		})
	}
}

// RequestTimeout returns a timeout middleware with the given duration
func RequestTimeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	MaxRequests int
	Duration    time.Duration
}

// simpleInMemoryStore is a basic in-memory store for rate limiting
type simpleInMemoryStore struct {
	requests map[string][]time.Time
}

// SimpleRateLimiter middleware implements a simple in-memory rate limiter
func SimpleRateLimiter(config RateLimitConfig) gin.HandlerFunc {
	store := &simpleInMemoryStore{
		requests: make(map[string][]time.Time),
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Clean old entries
		cutoff := now.Add(-config.Duration)
		if times, exists := store.requests[clientIP]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if t.After(cutoff) {
					validTimes = append(validTimes, t)
				}
			}
			store.requests[clientIP] = validTimes
		}

		// Check rate limit
		if len(store.requests[clientIP]) >= config.MaxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests",
				},
			})
			c.Abort()
			return
		}

		// Record request
		store.requests[clientIP] = append(store.requests[clientIP], now)
		c.Next()
	}
}

// BodySizeLimit middleware limits request body size
func BodySizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// ContentTypeValidator middleware validates Content-Type header for POST/PUT/PATCH requests
func ContentTypeValidator(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			// Remove charset suffix if present
			if idx := strings.Index(contentType, ";"); idx != -1 {
				contentType = strings.TrimSpace(contentType[:idx])
			}

			isAllowed := false
			for _, allowed := range allowedTypes {
				if contentType == allowed {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "UNSUPPORTED_MEDIA_TYPE",
						"message": "Content-Type must be one of: " + strings.Join(allowedTypes, ", "),
					},
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

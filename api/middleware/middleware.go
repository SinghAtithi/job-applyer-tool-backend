package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
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
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

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

// Timeout middleware adds request timeout
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if ctx.Err() == context.DeadlineExceeded {
			c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "TIMEOUT",
					"message": "Request timeout",
				},
			})
			return
		}
	}
}

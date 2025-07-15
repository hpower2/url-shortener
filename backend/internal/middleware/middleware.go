package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpower2/url-shortener/internal/errors"
	"github.com/hpower2/url-shortener/internal/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// Logger creates a structured logging middleware
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		entry := logger.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"body_size":  bodySize,
			"user_agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("Request completed")
		}
	}
}

// Recovery middleware with structured error handling
func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				logger.WithFields(logrus.Fields{
					"error":  err,
					"stack":  string(debug.Stack()),
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
				}).Error("Panic recovered")

				// Create structured error response
				appErr := errors.NewInternalError("Internal server error", fmt.Errorf("%v", err))
				c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS middleware with configurable origins
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RateLimiter creates a rate limiting middleware
func RateLimiter(rps float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			appErr := errors.NewRateLimitError("Rate limit exceeded", nil)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimiter creates a per-IP rate limiting middleware
func IPRateLimiter(rps float64, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			limiters[ip] = limiter
		}

		if !limiter.Allow() {
			appErr := errors.NewRateLimitError("Rate limit exceeded for IP", nil)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequestID middleware adds a unique request ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// Security middleware adds security headers
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

// ValidateContentType middleware validates content type for POST/PUT requests
func ValidateContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.Request.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				appErr := errors.NewBadRequestError("Content-Type must be application/json", nil)
				c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// ErrorHandler middleware handles application errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr := errors.GetAppError(err); appErr != nil {
				c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			} else {
				// Handle unknown errors
				appErr := errors.NewInternalError("Internal server error", err)
				c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			}
		}
	}
}

// Metrics middleware for collecting request metrics
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		// Here you would typically send metrics to your monitoring system
		// For now, we'll just set it in context for potential use
		c.Set("metrics", map[string]interface{}{
			"duration": duration,
			"status":   status,
			"method":   method,
			"path":     path,
		})
	}
}

// Timeout middleware adds request timeout
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Create a context with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// Replace the request context
		c.Request = c.Request.WithContext(timeoutCtx)

		// Channel to signal completion
		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// Request completed normally
			return
		case <-timeoutCtx.Done():
			// Request timed out
			appErr := errors.NewTimeoutError("Request timeout", timeoutCtx.Err())
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
		}
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())
}

// HealthCheck middleware for health check endpoints
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().Format(time.RFC3339),
				"version":   "1.0.0",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPWhitelist middleware allows only whitelisted IPs
func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		allowed := false
		for _, ip := range allowedIPs {
			if ip == clientIP {
				allowed = true
				break
			}
		}

		if !allowed {
			appErr := errors.NewForbiddenError("IP not allowed", nil)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}

		c.Next()
	}
}

// MaxBodySize middleware limits request body size
func MaxBodySize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			appErr := errors.NewBadRequestError("Request body too large", nil)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(authService interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := errors.NewUnauthorizedError("Authorization header required", nil)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			appErr := errors.NewUnauthorizedError("Invalid authorization header format", nil)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}

		token := authHeader[len(bearerPrefix):]
		if token == "" {
			appErr := errors.NewUnauthorizedError("Token is required", nil)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}

		// Cast authService to the correct type
		authSvc, ok := authService.(interface {
			ValidateToken(tokenString string) (*models.User, error)
		})
		if !ok {
			appErr := errors.NewInternalError("Invalid auth service configuration", nil)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}

		// Validate token
		user, err := authSvc.ValidateToken(token)
		if err != nil {
			appErr := errors.NewUnauthorizedError("Invalid token", err)
			c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Set("user", user)

		c.Next()
	}
}

// OptionalAuthMiddleware creates optional JWT authentication middleware
func OptionalAuthMiddleware(authService interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Extract token from "Bearer <token>" format
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.Next()
			return
		}

		token := authHeader[len(bearerPrefix):]
		if token == "" {
			c.Next()
			return
		}

		// Validate token (this will need to be updated when JWT is available)
		// user, err := authService.ValidateToken(token)
		// if err == nil {
		//     c.Set("user_id", user.ID)
		//     c.Set("user_email", user.Email)
		//     c.Set("user", user)
		// }

		c.Next()
	}
}

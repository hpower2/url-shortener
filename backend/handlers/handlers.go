package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpower2/url-shortener/internal/errors"
	"github.com/hpower2/url-shortener/internal/models"
	"github.com/hpower2/url-shortener/internal/services"
	"github.com/skip2/go-qrcode"
)

type Handler struct {
	urlService  services.URLService
	baseURL     string
	frontendURL string
}

func NewHandler(urlService services.URLService, baseURL, frontendURL string) *Handler {
	return &Handler{
		urlService:  urlService,
		baseURL:     baseURL,
		frontendURL: frontendURL,
	}
}

// CreateURL creates a new short URL with advanced features
func (h *Handler) CreateURL(c *gin.Context) {
	var req models.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get client info
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Create URL using service
	response, err := h.urlService.CreateURL(c.Request.Context(), &req, userID.(int), clientIP, userAgent)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RedirectURL redirects to original URL and records analytics
func (h *Handler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Get URL
	url, err := h.urlService.GetURL(c.Request.Context(), shortCode)
	if err != nil {
		h.ErrorPageHandler(c, err)
		return
	}

	// Record click with analytics
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referer := c.GetHeader("Referer")

	if err := h.urlService.RecordClick(c.Request.Context(), shortCode, clientIP, userAgent, referer); err != nil {
		// Log error but don't fail redirect
		// TODO: Add proper logging
	}

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

// GetURLStats returns detailed URL statistics
func (h *Handler) GetURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	stats, err := h.urlService.GetURLStats(c.Request.Context(), shortCode, userID.(int))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAllURLs returns paginated list of URLs
func (h *Handler) GetAllURLs(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	urls, total, err := h.urlService.GetAllURLs(c.Request.Context(), userID.(int), limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"urls":   urls,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateURL updates an existing URL
func (h *Handler) UpdateURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UpdateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url, err := h.urlService.UpdateURL(c.Request.Context(), shortCode, &req, userID.(int))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, url)
}

// DeleteURL deletes a URL
func (h *Handler) DeleteURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := h.urlService.DeleteURL(c.Request.Context(), shortCode, userID.(int)); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL deleted successfully"})
}

// GetAnalytics returns detailed analytics for a URL
func (h *Handler) GetAnalytics(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse days parameter
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
		return
	}

	analytics, err := h.urlService.GetAnalytics(c.Request.Context(), shortCode, userID.(int), days)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GenerateQRCode generates QR code for a URL
func (h *Handler) GenerateQRCode(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check ownership first
	url, err := h.urlService.GetURL(c.Request.Context(), shortCode)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Verify ownership (additional check)
	if url.UserID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Generate QR code for the short URL (not original URL)
	shortURL := fmt.Sprintf("%s/%s", h.baseURL, shortCode)

	// Generate QR code using the library
	qrCode, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Set appropriate headers for PNG image
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s-qr.png\"", shortCode))

	// Return the QR code image directly
	c.Data(http.StatusOK, "image/png", qrCode)
}

// HealthCheck returns service health status
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "2.0.0",
	})
}

// handleError handles different types of errors appropriately
func (h *Handler) handleError(c *gin.Context, err error) {
	if appErr := errors.GetAppError(err); appErr != nil {
		c.JSON(appErr.StatusCode, appErr.ToErrorResponse())
		return
	}

	// Fallback for unknown errors
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

// ErrorPageHandler handles errors for short URL redirects by redirecting to frontend
func (h *Handler) ErrorPageHandler(c *gin.Context, err error) {
	// Check if this is a short URL redirect request (not API)
	path := c.Request.URL.Path
	isShortURL := !strings.HasPrefix(path, "/api/") && !strings.HasPrefix(path, "/health")
	
	if !isShortURL {
		// For API requests, return JSON error
		h.handleError(c, err)
		return
	}
	
	// Extract short code from path
	shortCode := strings.TrimPrefix(path, "/")
	
	// For short URL requests, redirect to frontend error pages
	if appErr := errors.GetAppError(err); appErr != nil {
		switch appErr.Code {
		case errors.ErrCodeInactive:
			redirectURL := fmt.Sprintf("%s/error/inactive?code=%s", h.frontendURL, shortCode)
			c.Redirect(http.StatusFound, redirectURL)
		case errors.ErrCodeExpired:
			redirectURL := fmt.Sprintf("%s/error/expired?code=%s", h.frontendURL, shortCode)
			c.Redirect(http.StatusFound, redirectURL)
		case errors.ErrCodeNotFound:
			redirectURL := fmt.Sprintf("%s/error/not-found?code=%s", h.frontendURL, shortCode)
			c.Redirect(http.StatusFound, redirectURL)
		default:
			redirectURL := fmt.Sprintf("%s/error/server-error?code=%s", h.frontendURL, shortCode)
			c.Redirect(http.StatusFound, redirectURL)
		}
	} else {
		redirectURL := fmt.Sprintf("%s/error/server-error?code=%s", h.frontendURL, shortCode)
		c.Redirect(http.StatusFound, redirectURL)
	}
}

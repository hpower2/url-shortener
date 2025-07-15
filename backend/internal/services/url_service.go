package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hpower2/url-shortener/internal/errors"
	"github.com/hpower2/url-shortener/internal/models"
	"github.com/hpower2/url-shortener/internal/repository"
)

// URLService interface defines the contract for URL operations
type URLService interface {
	CreateURL(ctx context.Context, req *models.CreateURLRequest, userID int, clientIP, userAgent string) (*models.CreateURLResponse, error)
	GetURL(ctx context.Context, shortCode string) (*models.URL, error)
	GetURLStats(ctx context.Context, shortCode string, userID int) (*models.URLStatsResponse, error)
	GetAllURLs(ctx context.Context, userID int, limit, offset int) ([]models.URL, int, error)
	DeleteURL(ctx context.Context, shortCode string, userID int) error
	UpdateURL(ctx context.Context, shortCode string, req *models.UpdateURLRequest, userID int) (*models.URL, error)
	RecordClick(ctx context.Context, shortCode, clientIP, userAgent, referer string) error
	GetAnalytics(ctx context.Context, shortCode string, userID int, days int) (*models.URLAnalytics, error)
}

// urlService implements URLService interface
type urlService struct {
	urlRepo   repository.URLRepository
	userRepo  repository.UserRepository
	cacheRepo repository.CacheRepository
	baseURL   string
}

// NewURLService creates a new URL service
func NewURLService(urlRepo repository.URLRepository, userRepo repository.UserRepository, cacheRepo repository.CacheRepository, baseURL string) URLService {
	return &urlService{
		urlRepo:   urlRepo,
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		baseURL:   baseURL,
	}
}

// CreateURL creates a new short URL with user association
func (s *urlService) CreateURL(ctx context.Context, req *models.CreateURLRequest, userID int, clientIP, userAgent string) (*models.CreateURLResponse, error) {
	fmt.Println("Creating URL", req)
	fmt.Println("Client IP", clientIP)
	fmt.Println("User Agent", userAgent)
	fmt.Println("User ID", userID)
	fmt.Println("req.URL", req.URL)
	fmt.Println("req.CustomCode", req.CustomCode)
	fmt.Println("req.ExpiresAt", req.ExpiresAt)

	// Validate request
	if err := req.Validate(); err != nil {
		fmt.Println("Error", err)
		fmt.Println("Error type", reflect.TypeOf(err))

		return nil, errors.NewValidationError("Invalid request", err)
	}

	// Check if user can create more links
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get user", err)
	}

	if !user.CanCreateLink() {
		return nil, errors.NewValidationError(fmt.Sprintf("Link limit exceeded. You can create maximum %d links", user.LinkLimit), nil)
	}

	// Generate or use custom short code
	shortCode := req.CustomCode
	if shortCode == "" {
		var err error
		shortCode, err = s.generateUniqueShortCode(ctx)
		if err != nil {
			return nil, errors.NewInternalError("Failed to generate short code", err)
		}
	} else {
		// Check if custom code already exists
		exists, err := s.urlRepo.ExistsByShortCode(ctx, shortCode)
		if err != nil {
			return nil, errors.NewDatabaseError("Failed to check short code existence", err)
		}
		if exists {
			return nil, errors.NewAlreadyExistsError("Custom short code already exists", nil)
		}
	}

	// Create URL model
	url := &models.URL{
		ShortCode:   shortCode,
		OriginalURL: req.URL,
		UserID:      userID,
		IsActive:    true,
		ExpiresAt:   req.ExpiresAt.Time,
		IPAddress:   clientIP,
		UserAgent:   userAgent,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	createdURL, err := s.urlRepo.Create(ctx, url)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create URL", err)
	}

	// Cache the URL
	if err := s.cacheRepo.SetURL(ctx, shortCode, req.URL, 24*time.Hour); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache URL: %v\n", err)
	}

	// Create response
	response := &models.CreateURLResponse{
		ID:          createdURL.ID,
		ShortCode:   createdURL.ShortCode,
		OriginalURL: createdURL.OriginalURL,
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, createdURL.ShortCode),
		IsActive:    createdURL.IsActive,
		CreatedAt:   createdURL.CreatedAt,
		ExpiresAt:   createdURL.ExpiresAt,
		QRCode:      fmt.Sprintf("%s/api/v1/urls/%s/qr", s.baseURL, createdURL.ShortCode),
	}

	return response, nil
}

// GetURL retrieves a URL by short code
func (s *urlService) GetURL(ctx context.Context, shortCode string) (*models.URL, error) {
	if shortCode == "" {
		return nil, errors.NewValidationError("Short code is required", nil)
	}

	// Always get from database first to ensure we have the latest status
	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.NewNotFoundError("URL not found", err)
		}
		return nil, errors.NewDatabaseError("Failed to get URL", err)
	}

	// Check if URL is expired
	if url.IsExpired() {
		// Remove from cache if expired
		s.cacheRepo.DeleteURL(ctx, shortCode)
		return nil, errors.NewExpiredError("URL has expired", nil)
	}

	// Check if URL is active
	if !url.IsActive {
		// Remove from cache if inactive
		s.cacheRepo.DeleteURL(ctx, shortCode)
		return nil, errors.NewInactiveError("URL is not active", nil)
	}

	// Only cache if URL is active and not expired
	if err := s.cacheRepo.SetURL(ctx, shortCode, url.OriginalURL, 24*time.Hour); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache URL: %v\n", err)
	}

	return url, nil
}

// GetAllURLs retrieves all URLs with pagination
func (s *urlService) GetAllURLs(ctx context.Context, userID int, limit, offset int) ([]models.URL, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	urls, total, err := s.urlRepo.GetAllByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, errors.NewDatabaseError("Failed to get URLs", err)
	}

	return urls, total, nil
}

// DeleteURL deletes a URL by short code
func (s *urlService) DeleteURL(ctx context.Context, shortCode string, userID int) error {
	if shortCode == "" {
		return errors.NewValidationError("Short code is required", nil)
	}

	// Check ownership first
	owned, err := s.urlRepo.CheckOwnership(ctx, shortCode, userID)
	if err != nil {
		return errors.NewDatabaseError("Failed to check URL ownership", err)
	}
	if !owned {
		return errors.NewForbiddenError("URL not found or access denied", nil)
	}

	// Delete from cache first
	if err := s.cacheRepo.DeleteURL(ctx, shortCode); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to delete URL from cache: %v\n", err)
	}

	// Delete from database
	err = s.urlRepo.DeleteByUser(ctx, shortCode, userID)
	if err != nil {
		return errors.NewDatabaseError("Failed to delete URL", err)
	}

	return nil
}

// UpdateURL updates a URL
func (s *urlService) UpdateURL(ctx context.Context, shortCode string, req *models.UpdateURLRequest, userID int) (*models.URL, error) {
	if shortCode == "" {
		return nil, errors.NewValidationError("Short code is required", nil)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError("Invalid request", err)
	}

	// Check ownership first
	owned, err := s.urlRepo.CheckOwnership(ctx, shortCode, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to check URL ownership", err)
	}
	if !owned {
		return nil, errors.NewForbiddenError("URL not found or access denied", nil)
	}

	// Get existing URL
	url, err := s.urlRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get URL", err)
	}

	// Track if status-related fields are being changed
	statusChanged := false
	
	// Update fields
	if req.OriginalURL != "" {
		url.OriginalURL = req.OriginalURL
	}
	if req.IsActive != nil {
		if url.IsActive != *req.IsActive {
			statusChanged = true
		}
		url.IsActive = *req.IsActive
	}
	if req.ExpiresAt.Time != nil {
		if (url.ExpiresAt == nil && req.ExpiresAt.Time != nil) || 
		   (url.ExpiresAt != nil && req.ExpiresAt.Time != nil && !url.ExpiresAt.Equal(*req.ExpiresAt.Time)) {
			statusChanged = true
		}
		url.ExpiresAt = req.ExpiresAt.Time
	}
	url.UpdatedAt = time.Now()

	// Update in database
	updatedURL, err := s.urlRepo.Update(ctx, url)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to update URL", err)
	}

	// Clear cache if status changed or URL is inactive/expired
	if statusChanged || !updatedURL.IsActive || updatedURL.IsExpired() {
		if err := s.cacheRepo.DeleteURL(ctx, shortCode); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to delete URL from cache: %v\n", err)
		}
	} else {
		// Update cache only if URL is still active and not expired
		if err := s.cacheRepo.SetURL(ctx, shortCode, updatedURL.OriginalURL, 24*time.Hour); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to update URL in cache: %v\n", err)
		}
	}

	return updatedURL, nil
}

// RecordClick records a click event
func (s *urlService) RecordClick(ctx context.Context, shortCode, clientIP, userAgent, referer string) error {
	// Get URL
	url, err := s.GetURL(ctx, shortCode)
	if err != nil {
		return err
	}

	// Create click event
	clickEvent := &models.ClickEvent{
		URLId:     url.ID,
		IPAddress: clientIP,
		UserAgent: userAgent,
		Referer:   referer,
		ClickedAt: time.Now(),
	}

	// Save click event
	if err := s.urlRepo.CreateClickEvent(ctx, clickEvent); err != nil {
		return errors.NewDatabaseError("Failed to record click", err)
	}

	// Increment click count
	if err := s.urlRepo.IncrementClickCount(ctx, shortCode); err != nil {
		return errors.NewDatabaseError("Failed to increment click count", err)
	}

	// Increment click count in cache
	if err := s.cacheRepo.IncrementClickCount(ctx, shortCode); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to increment click count in cache: %v\n", err)
	}

	return nil
}

// GetURLStats retrieves URL statistics
func (s *urlService) GetURLStats(ctx context.Context, shortCode string, userID int) (*models.URLStatsResponse, error) {
	// Check ownership first
	owned, err := s.urlRepo.CheckOwnership(ctx, shortCode, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to check URL ownership", err)
	}
	if !owned {
		return nil, errors.NewForbiddenError("URL not found or access denied", nil)
	}

	// Get URL
	url, err := s.GetURL(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// Get analytics
	analytics, err := s.GetAnalytics(ctx, shortCode, userID, 30) // Get 30 days analytics
	if err != nil {
		return nil, err
	}

	// Get recent clicks
	recentClicks, err := s.urlRepo.GetClickEvents(ctx, url.ID, 10)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get recent clicks", err)
	}

	response := &models.URLStatsResponse{
		URL:          *url,
		TotalClicks:  analytics.TotalClicks,
		RecentClicks: recentClicks,
		Analytics:    *analytics,
	}

	return response, nil
}

// GetAnalytics retrieves URL analytics
func (s *urlService) GetAnalytics(ctx context.Context, shortCode string, userID int, days int) (*models.URLAnalytics, error) {
	// Check ownership first
	owned, err := s.urlRepo.CheckOwnership(ctx, shortCode, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to check URL ownership", err)
	}
	if !owned {
		return nil, errors.NewForbiddenError("URL not found or access denied", nil)
	}

	// Get URL
	url, err := s.GetURL(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// Get analytics data
	analytics, err := s.urlRepo.GetAnalyticsByUser(ctx, url.ID, userID, days)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to get analytics", err)
	}

	return analytics, nil
}

// generateUniqueShortCode generates a unique short code
func (s *urlService) generateUniqueShortCode(ctx context.Context) (string, error) {
	maxAttempts := 10

	for i := 0; i < maxAttempts; i++ {
		shortCode := s.generateShortCode()

		// Check if code already exists
		exists, err := s.urlRepo.ExistsByShortCode(ctx, shortCode)
		if err != nil {
			return "", err
		}

		if !exists {
			return shortCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique short code after %d attempts", maxAttempts)
}

// generateShortCode generates a random short code
func (s *urlService) generateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[s.randomInt(len(charset))]
	}

	return string(bytes)
}

// randomInt generates a random integer
func (s *urlService) randomInt(max int) int {
	bytes := make([]byte, 1)
	rand.Read(bytes)
	return int(bytes[0]) % max
}

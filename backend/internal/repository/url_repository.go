package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hpower2/url-shortener/database"
	"github.com/hpower2/url-shortener/internal/models"
)

// urlRepository implements URLRepository interface
type urlRepository struct {
	db *database.DB
}

// NewURLRepository creates a new URL repository
func NewURLRepository(db *database.DB) URLRepository {
	return &urlRepository{db: db}
}

// Create creates a new URL record
func (r *urlRepository) Create(ctx context.Context, url *models.URL) (*models.URL, error) {
	query := `
		INSERT INTO urls (short_code, original_url, user_id, is_active, expires_at, user_agent, ip_address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		url.ShortCode, url.OriginalURL, url.UserID, url.IsActive, url.ExpiresAt,
		url.UserAgent, url.IPAddress, url.CreatedAt, url.UpdatedAt,
	).Scan(&url.ID, &url.CreatedAt, &url.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}

	return url, nil
}

// GetByShortCode retrieves a URL by short code
func (r *urlRepository) GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	query := `
		SELECT id, short_code, original_url, user_id, created_at, updated_at, click_count, 
			   is_active, expires_at, user_agent, ip_address
		FROM urls 
		WHERE short_code = $1`

	url := &models.URL{}
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(
		&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID, &url.CreatedAt, &url.UpdatedAt,
		&url.ClickCount, &url.IsActive, &url.ExpiresAt, &url.UserAgent, &url.IPAddress,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("URL not found")
		}
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	return url, nil
}

// GetByID retrieves a URL by ID
func (r *urlRepository) GetByID(ctx context.Context, id int) (*models.URL, error) {
	query := `
		SELECT id, short_code, original_url, user_id, created_at, updated_at, click_count, 
			   is_active, expires_at, user_agent, ip_address
		FROM urls 
		WHERE id = $1`

	url := &models.URL{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID, &url.CreatedAt, &url.UpdatedAt,
		&url.ClickCount, &url.IsActive, &url.ExpiresAt, &url.UserAgent, &url.IPAddress,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("URL not found")
		}
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}

	return url, nil
}

// GetAll retrieves all URLs with pagination
func (r *urlRepository) GetAll(ctx context.Context, limit, offset int) ([]models.URL, int, error) {
	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM urls"
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get URLs with pagination
	query := `
		SELECT id, short_code, original_url, created_at, updated_at, click_count, 
			   is_active, expires_at, user_agent, ip_address
		FROM urls 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get URLs: %w", err)
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(
			&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.UpdatedAt,
			&url.ClickCount, &url.IsActive, &url.ExpiresAt, &url.UserAgent, &url.IPAddress,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}

	return urls, total, nil
}

// GetAllByUser retrieves all URLs for a specific user with pagination
func (r *urlRepository) GetAllByUser(ctx context.Context, userID int, limit, offset int) ([]models.URL, int, error) {
	// Get total count for the user
	var total int
	countQuery := `SELECT COUNT(*) FROM urls WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get URLs for the user
	query := `
		SELECT id, short_code, original_url, user_id, created_at, updated_at, click_count, 
			   is_active, expires_at, user_agent, ip_address
		FROM urls 
		WHERE user_id = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get URLs: %w", err)
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(
			&url.ID, &url.ShortCode, &url.OriginalURL, &url.UserID, &url.CreatedAt, &url.UpdatedAt,
			&url.ClickCount, &url.IsActive, &url.ExpiresAt, &url.UserAgent, &url.IPAddress,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan URL: %w", err)
		}
		urls = append(urls, url)
	}

	return urls, total, nil
}

// Update updates a URL record
func (r *urlRepository) Update(ctx context.Context, url *models.URL) (*models.URL, error) {
	query := `
		UPDATE urls 
		SET original_url = $2, is_active = $3, expires_at = $4, updated_at = $5
		WHERE short_code = $1
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		url.ShortCode, url.OriginalURL, url.IsActive, url.ExpiresAt, time.Now(),
	).Scan(&url.ID, &url.CreatedAt, &url.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to update URL: %w", err)
	}

	return url, nil
}

// Delete deletes a URL by short code
func (r *urlRepository) Delete(ctx context.Context, shortCode string) error {
	query := "DELETE FROM urls WHERE short_code = $1"
	result, err := r.db.ExecContext(ctx, query, shortCode)
	if err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("URL not found")
	}

	return nil
}

// DeleteByUser deletes a URL by short code for a specific user
func (r *urlRepository) DeleteByUser(ctx context.Context, shortCode string, userID int) error {
	query := `DELETE FROM urls WHERE short_code = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, shortCode, userID)
	if err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("URL not found or not owned by user")
	}

	return nil
}

// ExistsByShortCode checks if a URL exists by short code
func (r *urlRepository) ExistsByShortCode(ctx context.Context, shortCode string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)"
	var exists bool
	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check URL existence: %w", err)
	}
	return exists, nil
}

// IncrementClickCount increments the click count for a URL
func (r *urlRepository) IncrementClickCount(ctx context.Context, shortCode string) error {
	query := "UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1"
	_, err := r.db.ExecContext(ctx, query, shortCode)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}
	return nil
}

// CreateClickEvent creates a new click event record
func (r *urlRepository) CreateClickEvent(ctx context.Context, clickEvent *models.ClickEvent) error {
	query := `
		INSERT INTO click_events (url_id, ip_address, user_agent, referer, country, city, clicked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		clickEvent.URLId, clickEvent.IPAddress, clickEvent.UserAgent,
		clickEvent.Referer, clickEvent.Country, clickEvent.City, clickEvent.ClickedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create click event: %w", err)
	}

	return nil
}

// GetClickEvents retrieves click events for a URL
func (r *urlRepository) GetClickEvents(ctx context.Context, urlID int, limit int) ([]models.ClickEvent, error) {
	query := `
		SELECT id, url_id, ip_address, user_agent, referer, country, city, clicked_at
		FROM click_events 
		WHERE url_id = $1
		ORDER BY clicked_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, urlID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get click events: %w", err)
	}
	defer rows.Close()

	var events []models.ClickEvent
	for rows.Next() {
		var event models.ClickEvent
		err := rows.Scan(
			&event.ID, &event.URLId, &event.IPAddress, &event.UserAgent,
			&event.Referer, &event.Country, &event.City, &event.ClickedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan click event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

// GetAnalytics retrieves analytics data for a URL
func (r *urlRepository) GetAnalytics(ctx context.Context, urlID int, days int) (*models.URLAnalytics, error) {
	// For now, return basic analytics - you can enhance this with more complex queries
	analytics := &models.URLAnalytics{
		TotalClicks:    0,
		UniqueClicks:   0,
		ClicksToday:    0,
		ClicksThisWeek: 0,
		TopCountries:   []models.CountryStats{},
		TopReferrers:   []models.ReferrerStats{},
	}

	// Get total clicks
	query := "SELECT COUNT(*) FROM click_events WHERE url_id = $1"
	err := r.db.QueryRowContext(ctx, query, urlID).Scan(&analytics.TotalClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get unique clicks (unique IP addresses)
	query = "SELECT COUNT(DISTINCT ip_address) FROM click_events WHERE url_id = $1"
	err = r.db.QueryRowContext(ctx, query, urlID).Scan(&analytics.UniqueClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique clicks: %w", err)
	}

	// Get clicks today
	query = "SELECT COUNT(*) FROM click_events WHERE url_id = $1 AND clicked_at >= CURRENT_DATE"
	err = r.db.QueryRowContext(ctx, query, urlID).Scan(&analytics.ClicksToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks today: %w", err)
	}

	// Get clicks this week
	query = "SELECT COUNT(*) FROM click_events WHERE url_id = $1 AND clicked_at >= CURRENT_DATE - INTERVAL '7 days'"
	err = r.db.QueryRowContext(ctx, query, urlID).Scan(&analytics.ClicksThisWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks this week: %w", err)
	}

	return analytics, nil
}

// GetAnalyticsByUser retrieves URL analytics for a specific user
func (r *urlRepository) GetAnalyticsByUser(ctx context.Context, urlID int, userID int, days int) (*models.URLAnalytics, error) {
	// First check if the URL belongs to the user
	ownershipQuery := `SELECT COUNT(*) FROM urls WHERE id = $1 AND user_id = $2`
	var count int
	err := r.db.QueryRowContext(ctx, ownershipQuery, urlID, userID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check ownership: %w", err)
	}

	if count == 0 {
		return nil, fmt.Errorf("URL not found or not owned by user")
	}

	// Use the existing GetAnalytics method
	return r.GetAnalytics(ctx, urlID, days)
}

// CheckOwnership checks if a URL belongs to a specific user
func (r *urlRepository) CheckOwnership(ctx context.Context, shortCode string, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM urls WHERE short_code = $1 AND user_id = $2`
	var count int
	err := r.db.QueryRowContext(ctx, query, shortCode, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}

	return count > 0, nil
}

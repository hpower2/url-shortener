package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// OptionalTime is a custom type that can handle empty strings in JSON
type OptionalTime struct {
	*time.Time
}

// UnmarshalJSON implements json.Unmarshaler for OptionalTime
func (ot *OptionalTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	str := strings.Trim(string(data), "\"")

	// If empty string or null, set to nil
	if str == "" || str == "null" {
		ot.Time = nil
		return nil
	}

	// Try to parse the time
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	ot.Time = &t
	return nil
}

// MarshalJSON implements json.Marshaler for OptionalTime
func (ot OptionalTime) MarshalJSON() ([]byte, error) {
	if ot.Time == nil {
		return []byte("null"), nil
	}
	return json.Marshal(ot.Time.Format(time.RFC3339))
}

// URL represents a shortened URL record
type URL struct {
	ID          int        `db:"id" json:"id"`
	ShortCode   string     `db:"short_code" json:"short_code"`
	OriginalURL string     `db:"original_url" json:"original_url"`
	UserID      int        `db:"user_id" json:"user_id"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	ClickCount  int        `db:"click_count" json:"click_count"`
	IsActive    bool       `db:"is_active" json:"is_active"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	UserAgent   string     `db:"user_agent" json:"user_agent,omitempty"`
	IPAddress   string     `db:"ip_address" json:"ip_address,omitempty"`
}

// CreateURLRequest represents the request to create a new short URL
type CreateURLRequest struct {
	URL        string       `json:"url" binding:"required" validate:"required,url"`
	CustomCode string       `json:"custom_code,omitempty" validate:"omitempty,min=3,max=20,alphanum"`
	ExpiresAt  OptionalTime `json:"expires_at,omitempty"`
}

// CreateURLResponse represents the response when creating a short URL
type CreateURLResponse struct {
	ID          int        `json:"id"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	ShortURL    string     `json:"short_url"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	QRCode      string     `json:"qr_code_url,omitempty"`
}

// URLStatsResponse represents URL statistics
type URLStatsResponse struct {
	URL
	TotalClicks     int            `json:"total_clicks"`
	ClicksByCountry map[string]int `json:"clicks_by_country,omitempty"`
	ClicksByDate    map[string]int `json:"clicks_by_date,omitempty"`
	RecentClicks    []ClickEvent   `json:"recent_clicks,omitempty"`
	Analytics       URLAnalytics   `json:"analytics"`
}

// ClickEvent represents a click event
type ClickEvent struct {
	ID        int       `db:"id" json:"id"`
	URLId     int       `db:"url_id" json:"url_id"`
	IPAddress string    `db:"ip_address" json:"ip_address"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	Referer   string    `db:"referer" json:"referer"`
	Country   string    `db:"country" json:"country"`
	City      string    `db:"city" json:"city"`
	ClickedAt time.Time `db:"clicked_at" json:"clicked_at"`
}

// URLAnalytics represents analytics data
type URLAnalytics struct {
	TotalClicks    int             `json:"total_clicks"`
	UniqueClicks   int             `json:"unique_clicks"`
	ClicksToday    int             `json:"clicks_today"`
	ClicksThisWeek int             `json:"clicks_this_week"`
	TopCountries   []CountryStats  `json:"top_countries"`
	TopReferrers   []ReferrerStats `json:"top_referrers"`
}

// CountryStats represents click statistics by country
type CountryStats struct {
	Country string `json:"country"`
	Clicks  int    `json:"clicks"`
}

// ReferrerStats represents click statistics by referrer
type ReferrerStats struct {
	Referrer string `json:"referrer"`
	Clicks   int    `json:"clicks"`
}

// UpdateURLRequest represents the request to update a URL
type UpdateURLRequest struct {
	OriginalURL string       `json:"original_url,omitempty"`
	IsActive    *bool        `json:"is_active,omitempty"`
	ExpiresAt   OptionalTime `json:"expires_at,omitempty"`
}

// Validate validates the update URL request
func (req *UpdateURLRequest) Validate() error {
	if req.OriginalURL != "" {
		// Normalize URL
		req.OriginalURL = strings.TrimSpace(req.OriginalURL)
		if !strings.HasPrefix(req.OriginalURL, "http://") && !strings.HasPrefix(req.OriginalURL, "https://") {
			req.OriginalURL = "https://" + req.OriginalURL
		}

		// Validate URL format
		parsedURL, err := url.Parse(req.OriginalURL)
		if err != nil {
			return fmt.Errorf("invalid URL format: %w", err)
		}

		if parsedURL.Scheme == "" || parsedURL.Host == "" {
			return fmt.Errorf("URL must have scheme and host")
		}
	}

	// Validate expiration date
	if req.ExpiresAt.Time != nil && req.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("expiration date cannot be in the past")
	}

	return nil
}

// Validate validates the URL model
func (u *URL) Validate() error {
	if u.ShortCode == "" {
		return fmt.Errorf("short code is required")
	}

	if u.OriginalURL == "" {
		return fmt.Errorf("original URL is required")
	}

	if !u.IsValidURL() {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// IsValidURL checks if the original URL is valid
func (u *URL) IsValidURL() bool {
	parsedURL, err := url.Parse(u.OriginalURL)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

// IsExpired checks if the URL has expired
func (u *URL) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}

// NormalizeURL normalizes the original URL
func (u *URL) NormalizeURL() {
	u.OriginalURL = strings.TrimSpace(u.OriginalURL)
	if !strings.HasPrefix(u.OriginalURL, "http://") && !strings.HasPrefix(u.OriginalURL, "https://") {
		u.OriginalURL = "https://" + u.OriginalURL
	}
}

// Validate validates the create URL request
func (req *CreateURLRequest) Validate() error {
	if req.URL == "" {
		return fmt.Errorf("URL is required")
	}

	// Normalize URL
	req.URL = strings.TrimSpace(req.URL)
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		req.URL = "https://" + req.URL
	}

	// Validate URL format
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("URL must have scheme and host")
	}

	// Validate custom code if provided
	if req.CustomCode != "" {
		if len(req.CustomCode) < 3 || len(req.CustomCode) > 20 {
			return fmt.Errorf("custom code must be between 3 and 20 characters")
		}

		// Check if custom code contains only alphanumeric characters
		for _, char := range req.CustomCode {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
				return fmt.Errorf("custom code must contain only alphanumeric characters")
			}
		}
	}

	// Validate expiration date
	if req.ExpiresAt.Time != nil && req.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("expiration date cannot be in the past")
	}

	return nil
}

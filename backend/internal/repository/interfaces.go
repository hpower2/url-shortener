package repository

import (
	"context"
	"time"

	"github.com/hpower2/url-shortener/internal/models"
)

// URLRepository interface defines the contract for URL database operations
type URLRepository interface {
	Create(ctx context.Context, url *models.URL) (*models.URL, error)
	GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error)
	GetByID(ctx context.Context, id int) (*models.URL, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.URL, int, error)
	GetAllByUser(ctx context.Context, userID int, limit, offset int) ([]models.URL, int, error)
	Update(ctx context.Context, url *models.URL) (*models.URL, error)
	Delete(ctx context.Context, shortCode string) error
	DeleteByUser(ctx context.Context, shortCode string, userID int) error
	ExistsByShortCode(ctx context.Context, shortCode string) (bool, error)
	IncrementClickCount(ctx context.Context, shortCode string) error
	CreateClickEvent(ctx context.Context, clickEvent *models.ClickEvent) error
	GetClickEvents(ctx context.Context, urlID int, limit int) ([]models.ClickEvent, error)
	GetAnalytics(ctx context.Context, urlID int, days int) (*models.URLAnalytics, error)
	GetAnalyticsByUser(ctx context.Context, urlID int, userID int, days int) (*models.URLAnalytics, error)
	CheckOwnership(ctx context.Context, shortCode string, userID int) (bool, error)
}

// CacheRepository interface defines the contract for cache operations
type CacheRepository interface {
	SetURL(ctx context.Context, shortCode, originalURL string, expiration time.Duration) error
	GetURL(ctx context.Context, shortCode string) (string, error)
	DeleteURL(ctx context.Context, shortCode string) error
	IncrementClickCount(ctx context.Context, shortCode string) error
	GetClickCount(ctx context.Context, shortCode string) (int64, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
} 
package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/hpower2/url-shortener/config"
)

type DB struct {
	*sqlx.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// Example URL shortener models and queries
type URL struct {
	ID          int    `db:"id" json:"id"`
	ShortCode   string `db:"short_code" json:"short_code"`
	OriginalURL string `db:"original_url" json:"original_url"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	ClickCount  int    `db:"click_count" json:"click_count"`
}

func (db *DB) CreateURL(shortCode, originalURL string) (*URL, error) {
	query := `
		INSERT INTO urls (short_code, original_url) 
		VALUES ($1, $2) 
		RETURNING id, short_code, original_url, created_at, click_count
	`
	
	var url URL
	err := db.Get(&url, query, shortCode, originalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL: %w", err)
	}
	
	return &url, nil
}

func (db *DB) GetURLByShortCode(shortCode string) (*URL, error) {
	query := `SELECT id, short_code, original_url, created_at, click_count FROM urls WHERE short_code = $1`
	
	var url URL
	err := db.Get(&url, query, shortCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	
	return &url, nil
}

func (db *DB) IncrementClickCount(shortCode string) error {
	query := `UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1`
	
	_, err := db.Exec(query, shortCode)
	if err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}
	
	return nil
}

func (db *DB) GetAllURLs() ([]URL, error) {
	query := `SELECT id, short_code, original_url, created_at, click_count FROM urls ORDER BY created_at DESC`
	
	var urls []URL
	err := db.Select(&urls, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all URLs: %w", err)
	}
	
	return urls, nil
} 
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Security SecurityConfig `json:"security"`
	Logging  LoggingConfig  `json:"logging"`
	App      AppConfig      `json:"app"`
	SMTP     SMTPConfig     `json:"smtp"`
	RabbitMQ RabbitMQConfig `json:"rabbitmq"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port            string        `json:"port"`
	Host            string        `json:"host"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	MaxHeaderBytes  int           `json:"max_header_bytes"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            string        `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	DBName          string        `json:"db_name"`
	SSLMode         string        `json:"ssl_mode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	Password     string        `json:"password"`
	DB           int           `json:"db"`
	PoolSize     int           `json:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	JWTSecret      string        `json:"jwt_secret"`
	JWTExpiration  time.Duration `json:"jwt_expiration"`
	RateLimitRPS   float64       `json:"rate_limit_rps"`
	RateLimitBurst int           `json:"rate_limit_burst"`
	MaxRequestSize int64         `json:"max_request_size"`
	AllowedOrigins []string      `json:"allowed_origins"`
	TrustedProxies []string      `json:"trusted_proxies"`
	EnableHTTPS    bool          `json:"enable_https"`
	CertFile       string        `json:"cert_file"`
	KeyFile        string        `json:"key_file"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// AppConfig represents application-specific configuration
type AppConfig struct {
	Name                string        `json:"name"`
	Version             string        `json:"version"`
	Environment         string        `json:"environment"`
	BaseURL             string        `json:"base_url"`
	FrontendURL         string        `json:"frontend_url"`
	ShortCodeLength     int           `json:"short_code_length"`
	DefaultExpiration   time.Duration `json:"default_expiration"`
	MaxCustomCodeLength int           `json:"max_custom_code_length"`
	EnableAnalytics     bool          `json:"enable_analytics"`
	EnableQRCode        bool          `json:"enable_qr_code"`
	CleanupInterval     time.Duration `json:"cleanup_interval"`
}

// SMTPConfig represents SMTP configuration
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

// RabbitMQConfig represents RabbitMQ configuration
type RabbitMQConfig struct {
	URL      string `json:"url"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:     getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
			MaxHeaderBytes:  getIntEnv("SERVER_MAX_HEADER_BYTES", 1<<20), // 1MB
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "password"),
			DBName:          getEnv("DB_NAME", "urlshortener"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnv("REDIS_PORT", "6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getIntEnv("REDIS_DB", 0),
			PoolSize:     getIntEnv("REDIS_POOL_SIZE", 10),
			MinIdleConns: getIntEnv("REDIS_MIN_IDLE_CONNS", 5),
			DialTimeout:  getDurationEnv("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  getDurationEnv("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getDurationEnv("REDIS_WRITE_TIMEOUT", 3*time.Second),
		},
		Security: SecurityConfig{
			JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
			JWTExpiration:  getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
			RateLimitRPS:   getFloat64Env("RATE_LIMIT_RPS", 10.0),
			RateLimitBurst: getIntEnv("RATE_LIMIT_BURST", 20),
			MaxRequestSize: getInt64Env("MAX_REQUEST_SIZE", 1<<20), // 1MB
			AllowedOrigins: getSliceEnv("ALLOWED_ORIGINS", []string{"*"}),
			TrustedProxies: getSliceEnv("TRUSTED_PROXIES", []string{}),
			EnableHTTPS:    getBoolEnv("ENABLE_HTTPS", false),
			CertFile:       getEnv("CERT_FILE", ""),
			KeyFile:        getEnv("KEY_FILE", ""),
		},
		Logging: LoggingConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Format:     getEnv("LOG_FORMAT", "json"),
			Output:     getEnv("LOG_OUTPUT", "stdout"),
			MaxSize:    getIntEnv("LOG_MAX_SIZE", 100),
			MaxBackups: getIntEnv("LOG_MAX_BACKUPS", 3),
			MaxAge:     getIntEnv("LOG_MAX_AGE", 28),
			Compress:   getBoolEnv("LOG_COMPRESS", true),
		},
		App: AppConfig{
			Name:                getEnv("APP_NAME", "URL Shortener"),
			Version:             getEnv("APP_VERSION", "1.0.0"),
			Environment:         getEnv("APP_ENV", "development"),
			BaseURL:             getEnv("BASE_URL", "http://localhost:8080"),
			FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
			ShortCodeLength:     getIntEnv("SHORT_CODE_LENGTH", 8),
			DefaultExpiration:   getDurationEnv("DEFAULT_EXPIRATION", 0), // 0 means no expiration
			MaxCustomCodeLength: getIntEnv("MAX_CUSTOM_CODE_LENGTH", 20),
			EnableAnalytics:     getBoolEnv("ENABLE_ANALYTICS", true),
			EnableQRCode:        getBoolEnv("ENABLE_QR_CODE", true),
			CleanupInterval:     getDurationEnv("CLEANUP_INTERVAL", 24*time.Hour),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.hostinger.com"),
			Port:     getIntEnv("SMTP_PORT", 465),
			Username: getEnv("SMTP_USERNAME", "me@irvineafri.com"),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@irvineafri.com"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:      getEnv("RABBITMQ_URL", ""),
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			Username: getEnv("RABBITMQ_USERNAME", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
		},
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validate database config
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}

	// Validate Redis config
	if c.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}

	// Validate security config
	if c.Security.JWTSecret == "" || c.Security.JWTSecret == "your-secret-key" {
		return fmt.Errorf("JWT secret must be set and not be default value")
	}

	// Validate app config
	if c.App.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	if c.App.ShortCodeLength < 4 || c.App.ShortCodeLength > 20 {
		return fmt.Errorf("short code length must be between 4 and 20")
	}

	return nil
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User,
		c.Database.Password, c.Database.DBName, c.Database.SSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getFloat64Env(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated parsing
		var result []string
		for _, item := range strings.Split(value, ",") {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

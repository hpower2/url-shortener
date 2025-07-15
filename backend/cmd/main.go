package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	oldconfig "github.com/hpower2/url-shortener/config"
	"github.com/hpower2/url-shortener/database"
	"github.com/hpower2/url-shortener/handlers"
	"github.com/hpower2/url-shortener/internal/config"
	"github.com/hpower2/url-shortener/internal/middleware"
	"github.com/hpower2/url-shortener/internal/repository"
	"github.com/hpower2/url-shortener/internal/services"
	"github.com/hpower2/url-shortener/redis"
	"github.com/sirupsen/logrus"
)

// convertDatabaseConfig converts new config to old config format
func convertDatabaseConfig(newCfg *config.DatabaseConfig) *oldconfig.DatabaseConfig {
	return &oldconfig.DatabaseConfig{
		Host:     newCfg.Host,
		Port:     newCfg.Port,
		User:     newCfg.User,
		Password: newCfg.Password,
		DBName:   newCfg.DBName,
		SSLMode:  newCfg.SSLMode,
	}
}

// convertRedisConfig converts new config to old config format
func convertRedisConfig(newCfg *config.RedisConfig) *oldconfig.RedisConfig {
	return &oldconfig.RedisConfig{
		Host:     newCfg.Host,
		Port:     newCfg.Port,
		Password: newCfg.Password,
		DB:       newCfg.DB,
	}
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Initialize database
	db, err := database.NewDatabase(convertDatabaseConfig(&cfg.Database))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := redis.NewRedisClient(convertRedisConfig(&cfg.Redis))
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	urlRepo := repository.NewURLRepository(db)
	cacheRepo := repository.NewCacheRepository(redisClient)
	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)

	// Initialize services
	baseURL := cfg.App.BaseURL
	urlService := services.NewURLService(urlRepo, userRepo, cacheRepo, baseURL)
	authService := services.NewAuthService(userRepo, cfg.Security.JWTSecret)
	emailService := services.NewEmailService(&cfg.SMTP)
	otpService := services.NewOTPService(otpRepo, userRepo)
	rabbitMQService := services.NewRabbitMQService(&cfg.RabbitMQ)
	emailQueueConsumer := services.NewEmailQueueConsumer(rabbitMQService, emailService, otpService, cfg)

	// Initialize handlers
	handler := handlers.NewHandler(urlService, baseURL, cfg.App.FrontendURL)
	authHandler := handlers.NewAuthHandler(authService)
	otpHandler := handlers.NewOTPHandler(otpService, emailQueueConsumer, userRepo)

	// Start email queue consumer
	ctx := context.Background()
	if err := emailQueueConsumer.Start(ctx); err != nil {
		log.Printf("Failed to start email queue consumer: %v", err)
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS([]string{"*"}))
	router.Use(middleware.RateLimiter(100, 10)) // 100 requests per second, burst of 10
	router.Use(middleware.RequestID())
	router.Use(middleware.Security())

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		// OTP routes (public)
		otp := api.Group("/otp")
		{
			otp.POST("/generate", otpHandler.GenerateOTP)
			otp.POST("/verify", otpHandler.VerifyOTP)
		}

		// Protected routes (require authentication)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User profile routes
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", authHandler.UpdateProfile)
			protected.POST("/profile/change-password", authHandler.ChangePassword)
			protected.POST("/auth/refresh", authHandler.RefreshToken)

			// URL management (protected)
			protected.POST("/urls", handler.CreateURL)
			protected.GET("/urls", handler.GetAllURLs)
			protected.GET("/urls/:shortCode", handler.GetURLStats)
			protected.PUT("/urls/:shortCode", handler.UpdateURL)
			protected.DELETE("/urls/:shortCode", handler.DeleteURL)

			// Analytics (protected)
			protected.GET("/urls/:shortCode/analytics", handler.GetAnalytics)

			// QR Code generation (protected)
			protected.GET("/urls/:shortCode/qr", handler.GenerateQRCode)
		}
	}

	// Direct redirect routes (must be last to avoid conflicts and remain public)
	router.GET("/:shortCode", handler.RedirectURL)

	// Start server
	log.Printf("🚀 URL Shortener v2.0 starting on port %s", cfg.Server.Port)
	log.Printf("📊 Features enabled: Custom codes, Analytics, QR codes, Rate limiting, User Authentication")
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

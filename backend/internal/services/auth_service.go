package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hpower2/url-shortener/internal/errors"
	"github.com/hpower2/url-shortener/internal/models"
	"github.com/hpower2/url-shortener/internal/repository"
)

// AuthService interface defines the contract for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	ValidateToken(tokenString string) (*models.User, error)
	RefreshToken(ctx context.Context, userID int) (string, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	UpdateUser(ctx context.Context, userID int, req *models.UpdateUserRequest) (*models.User, error)
	ChangePassword(ctx context.Context, userID int, req *models.ChangePasswordRequest) error
}

// authService implements AuthService interface
type authService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.LoginResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		log.Println("Invalid registration data", err)
		return nil, errors.NewValidationError("Invalid registration data", err)
	}

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		log.Println("Failed to check user existence", err)
		return nil, errors.NewDatabaseError("Failed to check user existence", err)
	}
	if exists {
		log.Println("User with this email already exists")
		return nil, errors.NewAlreadyExistsError("User with this email already exists", nil)
	}

	// Create user
	user := &models.User{
		Email:         req.Email,
		Password:      req.Password,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		IsActive:      true,
		EmailVerified: false, // User needs to verify email first
		LinkCount:     0,
		LinkLimit:     50,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, errors.NewInternalError("Failed to hash password", err)
	}

	// Save user
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to create user", err)
	}

	// Generate JWT token
	token, err := s.generateToken(createdUser)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate token", err)
	}

	return &models.LoginResponse{
		User:  createdUser.ToResponse(),
		Token: token,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError("Invalid login data", err)
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewNotFoundError("Invalid email or password", nil)
	}

	// Check if user is active
	if !user.IsValidForLogin() {
		return nil, errors.NewUnauthorizedError("Account is deactivated", nil)
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		return nil, errors.NewUnauthorizedError("Invalid email or password", nil)
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate token", err)
	}

	return &models.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *authService) ValidateToken(tokenString string) (*models.User, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid token", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.NewUnauthorizedError("Invalid token claims", nil)
	}

	// Get user from database
	user, err := s.userRepo.GetByID(context.Background(), claims.UserID)
	if err != nil {
		return nil, errors.NewUnauthorizedError("User not found", err)
	}

	// Check if user is still active
	if !user.IsValidForLogin() {
		return nil, errors.NewUnauthorizedError("Account is deactivated", nil)
	}

	return user, nil
}

// RefreshToken generates a new JWT token for a user
func (s *authService) RefreshToken(ctx context.Context, userID int) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", errors.NewNotFoundError("User not found", err)
	}

	if !user.IsValidForLogin() {
		return "", errors.NewUnauthorizedError("Account is deactivated", nil)
	}

	return s.generateToken(user)
}

// GetUserByID retrieves a user by ID
func (s *authService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("User not found", err)
	}
	return user, nil
}

// UpdateUser updates user information
func (s *authService) UpdateUser(ctx context.Context, userID int, req *models.UpdateUserRequest) (*models.User, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError("Invalid update data", err)
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFoundError("User not found", err)
	}

	// Update fields
	if req.Email != "" {
		// Check if email is already taken by another user
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && existingUser.ID != userID {
			return nil, errors.NewAlreadyExistsError("Email already in use", nil)
		}
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	user.UpdatedAt = time.Now()

	// Update user
	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to update user", err)
	}

	return updatedUser, nil
}

// ChangePassword changes user password
func (s *authService) ChangePassword(ctx context.Context, userID int, req *models.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.NewNotFoundError("User not found", err)
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		return errors.NewUnauthorizedError("Current password is incorrect", nil)
	}

	// Validate new password
	if len(req.NewPassword) < 8 {
		return errors.NewValidationError("New password must be at least 8 characters long", nil)
	}

	// Hash new password
	user.Password = req.NewPassword
	if err := user.HashPassword(); err != nil {
		return errors.NewInternalError("Failed to hash password", err)
	}

	// Update user
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return errors.NewDatabaseError("Failed to update password", err)
	}

	return nil
}

// generateToken generates a JWT token for a user
func (s *authService) generateToken(user *models.User) (string, error) {
	// Create claims
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "url-shortener",
			Subject:   fmt.Sprintf("user-%d", user.ID),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

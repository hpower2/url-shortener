package models

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID              int        `db:"id" json:"id"`
	Email           string     `db:"email" json:"email"`
	Password        string     `db:"password" json:"-"` // Never expose password in JSON
	FirstName       string     `db:"first_name" json:"first_name"`
	LastName        string     `db:"last_name" json:"last_name"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	EmailVerified   bool       `db:"email_verified" json:"email_verified"`
	EmailVerifiedAt *time.Time `db:"email_verified_at" json:"email_verified_at,omitempty"`
	LinkCount       int        `db:"link_count" json:"link_count"`
	LinkLimit       int        `db:"link_limit" json:"link_limit"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required" validate:"required,email"`
	Password  string `json:"password" binding:"required" validate:"required,min=8"`
	FirstName string `json:"first_name" binding:"required" validate:"required,min=2"`
	LastName  string `json:"last_name" binding:"required" validate:"required,min=2"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserResponse represents user data in responses (without sensitive info)
type UserResponse struct {
	ID              int        `json:"id"`
	Email           string     `json:"email"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	IsActive        bool       `json:"is_active"`
	EmailVerified   bool       `json:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LinkCount       int        `json:"link_count"`
	LinkLimit       int        `json:"link_limit"`
	CreatedAt       time.Time  `json:"created_at"`
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty,min=2"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=2"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8"`
}

// Validate validates the registration request
func (req *RegisterRequest) Validate() error {
	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Validate email format
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("invalid email format")
	}

	// Validate password strength
	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Validate names
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if len(req.FirstName) < 2 {
		return fmt.Errorf("first name must be at least 2 characters long")
	}

	if len(req.LastName) < 2 {
		return fmt.Errorf("last name must be at least 2 characters long")
	}

	return nil
}

// Validate validates the login request
func (req *LoginRequest) Validate() error {
	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Validate email format
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fmt.Errorf("invalid email format")
	}

	if req.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

// Validate validates the update user request
func (req *UpdateUserRequest) Validate() error {
	if req.Email != "" {
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return fmt.Errorf("invalid email format")
		}
	}

	if req.FirstName != "" {
		req.FirstName = strings.TrimSpace(req.FirstName)
		if len(req.FirstName) < 2 {
			return fmt.Errorf("first name must be at least 2 characters long")
		}
	}

	if req.LastName != "" {
		req.LastName = strings.TrimSpace(req.LastName)
		if len(req.LastName) < 2 {
			return fmt.Errorf("last name must be at least 2 characters long")
		}
	}

	return nil
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// IsValidForLogin checks if user can login
// For initial implementation, allow login without email verification
// This will be updated later to require email verification
func (u *User) IsValidForLogin() bool {
	return u.IsActive
}

// CanCreateLink checks if user can create more links
func (u *User) CanCreateLink() bool {
	return u.LinkCount < u.LinkLimit
}

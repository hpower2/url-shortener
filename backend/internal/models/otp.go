package models

import (
	"fmt"
	"time"
)

// OTPVerification represents an OTP verification record
type OTPVerification struct {
	ID         int        `db:"id" json:"id"`
	UserID     int        `db:"user_id" json:"user_id"`
	Email      string     `db:"email" json:"email"`
	OTPCode    string     `db:"otp_code" json:"otp_code"`
	Purpose    string     `db:"purpose" json:"purpose"`
	IsVerified bool       `db:"is_verified" json:"is_verified"`
	ExpiresAt  time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	VerifiedAt *time.Time `db:"verified_at" json:"verified_at,omitempty"`
}

// OTPRequest represents a request to generate OTP
type OTPRequest struct {
	Email   string `json:"email" binding:"required" validate:"required,email"`
	Purpose string `json:"purpose" binding:"required" validate:"required"`
}

// OTPVerifyRequest represents a request to verify OTP
type OTPVerifyRequest struct {
	Email   string `json:"email" binding:"required" validate:"required,email"`
	OTPCode string `json:"otp_code" binding:"required" validate:"required,len=6"`
	Purpose string `json:"purpose" binding:"required" validate:"required"`
}

// OTPResponse represents the response after OTP generation
type OTPResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

// OTPVerifyResponse represents the response after OTP verification
type OTPVerifyResponse struct {
	Message    string `json:"message"`
	IsVerified bool   `json:"is_verified"`
}

// Validate validates the OTP request
func (req *OTPRequest) Validate() error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Purpose == "" {
		return fmt.Errorf("purpose is required")
	}
	if req.Purpose != "email_verification" && req.Purpose != "password_reset" {
		return fmt.Errorf("invalid purpose")
	}
	return nil
}

// Validate validates the OTP verify request
func (req *OTPVerifyRequest) Validate() error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.OTPCode == "" {
		return fmt.Errorf("OTP code is required")
	}
	if len(req.OTPCode) != 6 {
		return fmt.Errorf("OTP code must be 6 digits")
	}
	if req.Purpose == "" {
		return fmt.Errorf("purpose is required")
	}
	return nil
}

// IsExpired checks if the OTP has expired
func (otp *OTPVerification) IsExpired() bool {
	return time.Now().After(otp.ExpiresAt)
}

// CanBeVerified checks if the OTP can be verified
func (otp *OTPVerification) CanBeVerified() bool {
	return !otp.IsExpired() && !otp.IsVerified
} 
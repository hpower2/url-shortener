package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/hpower2/url-shortener/internal/errors"
	"github.com/hpower2/url-shortener/internal/models"
	"github.com/hpower2/url-shortener/internal/repository"
)

// OTPService interface defines the contract for OTP operations
type OTPService interface {
	GenerateOTP(ctx context.Context, userID int, email, purpose string) (*models.OTPResponse, error)
	VerifyOTP(ctx context.Context, req *models.OTPVerifyRequest) (*models.OTPVerifyResponse, error)
	CleanupExpiredOTPs(ctx context.Context) error
}

// otpService implements OTPService interface
type otpService struct {
	otpRepo  repository.OTPRepository
	userRepo repository.UserRepository
}

// NewOTPService creates a new OTP service
func NewOTPService(otpRepo repository.OTPRepository, userRepo repository.UserRepository) OTPService {
	return &otpService{
		otpRepo:  otpRepo,
		userRepo: userRepo,
	}
}

// GenerateOTP generates a new OTP for the user
func (s *otpService) GenerateOTP(ctx context.Context, userID int, email, purpose string) (*models.OTPResponse, error) {
	// Generate 6-digit OTP
	otpCode, err := s.generateOTPCode()
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate OTP code", err)
	}

	// Set expiration time (10 minutes from now)
	expiresAt := time.Now().Add(10 * time.Minute)

	// Create OTP record
	otp := &models.OTPVerification{
		UserID:    userID,
		Email:     email,
		OTPCode:   otpCode,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	// Save OTP to database (this will replace any existing OTP for the same user/purpose)
	createdOTP, err := s.otpRepo.Create(ctx, otp)
	if err != nil {
		return nil, errors.NewDatabaseError("Failed to save OTP", err)
	}

	return &models.OTPResponse{
		Message:   "OTP sent successfully",
		ExpiresAt: createdOTP.ExpiresAt,
	}, nil
}

// VerifyOTP verifies the provided OTP
func (s *otpService) VerifyOTP(ctx context.Context, req *models.OTPVerifyRequest) (*models.OTPVerifyResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, errors.NewValidationError("Invalid OTP verification request", err)
	}

	// Get OTP record
	otp, err := s.otpRepo.GetByEmailAndPurpose(ctx, req.Email, req.Purpose)
	if err != nil {
		return &models.OTPVerifyResponse{
			Message:    "Invalid or expired OTP",
			IsVerified: false,
		}, nil
	}

	// Check if OTP can be verified
	if !otp.CanBeVerified() {
		return &models.OTPVerifyResponse{
			Message:    "OTP has expired or already been used",
			IsVerified: false,
		}, nil
	}

	// Verify OTP code
	if otp.OTPCode != req.OTPCode {
		return &models.OTPVerifyResponse{
			Message:    "Invalid OTP code",
			IsVerified: false,
		}, nil
	}

	// Mark OTP as verified
	now := time.Now()
	otp.IsVerified = true
	otp.VerifiedAt = &now

	if err := s.otpRepo.Update(ctx, otp); err != nil {
		return nil, errors.NewDatabaseError("Failed to update OTP status", err)
	}

	// If this is email verification, update user's email verification status
	if req.Purpose == "email_verification" {
		user, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return nil, errors.NewDatabaseError("Failed to get user", err)
		}

		user.EmailVerified = true
		user.EmailVerifiedAt = &now

		if _, err := s.userRepo.Update(ctx, user); err != nil {
			return nil, errors.NewDatabaseError("Failed to update user verification status", err)
		}
	}

	return &models.OTPVerifyResponse{
		Message:    "OTP verified successfully",
		IsVerified: true,
	}, nil
}

// CleanupExpiredOTPs removes expired OTP records
func (s *otpService) CleanupExpiredOTPs(ctx context.Context) error {
	return s.otpRepo.DeleteExpired(ctx)
}

// generateOTPCode generates a 6-digit OTP code
func (s *otpService) generateOTPCode() (string, error) {
	// Generate random 6-digit number
	max := big.NewInt(999999)
	min := big.NewInt(100000)

	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Add(n, min).Int64()), nil
}

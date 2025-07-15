package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hpower2/url-shortener/database"
	"github.com/hpower2/url-shortener/internal/models"
)

// OTPRepository interface defines the contract for OTP database operations
type OTPRepository interface {
	Create(ctx context.Context, otp *models.OTPVerification) (*models.OTPVerification, error)
	GetByEmailAndPurpose(ctx context.Context, email, purpose string) (*models.OTPVerification, error)
	Update(ctx context.Context, otp *models.OTPVerification) error
	DeleteExpired(ctx context.Context) error
	DeleteByUserAndPurpose(ctx context.Context, userID int, purpose string) error
}

// otpRepository implements OTPRepository interface
type otpRepository struct {
	db *database.DB
}

// NewOTPRepository creates a new OTP repository
func NewOTPRepository(db *database.DB) OTPRepository {
	return &otpRepository{db: db}
}

// Create creates a new OTP record, replacing any existing unverified OTP for the same user/purpose
func (r *otpRepository) Create(ctx context.Context, otp *models.OTPVerification) (*models.OTPVerification, error) {
	// First, delete any existing unverified OTP for this user/purpose
	deleteQuery := `
		DELETE FROM otp_verifications 
		WHERE user_id = $1 AND purpose = $2 AND is_verified = FALSE`
	
	_, err := r.db.ExecContext(ctx, deleteQuery, otp.UserID, otp.Purpose)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing OTP: %w", err)
	}

	// Create new OTP record
	query := `
		INSERT INTO otp_verifications (user_id, email, otp_code, purpose, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err = r.db.QueryRowContext(ctx, query,
		otp.UserID, otp.Email, otp.OTPCode, otp.Purpose, otp.ExpiresAt, otp.CreatedAt,
	).Scan(&otp.ID, &otp.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create OTP: %w", err)
	}

	return otp, nil
}

// GetByEmailAndPurpose retrieves the latest unverified OTP for email and purpose
func (r *otpRepository) GetByEmailAndPurpose(ctx context.Context, email, purpose string) (*models.OTPVerification, error) {
	query := `
		SELECT id, user_id, email, otp_code, purpose, is_verified, expires_at, created_at, verified_at
		FROM otp_verifications 
		WHERE email = $1 AND purpose = $2 AND is_verified = FALSE
		ORDER BY created_at DESC
		LIMIT 1`

	otp := &models.OTPVerification{}
	err := r.db.QueryRowContext(ctx, query, email, purpose).Scan(
		&otp.ID, &otp.UserID, &otp.Email, &otp.OTPCode, &otp.Purpose, 
		&otp.IsVerified, &otp.ExpiresAt, &otp.CreatedAt, &otp.VerifiedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("OTP not found")
		}
		return nil, fmt.Errorf("failed to get OTP: %w", err)
	}

	return otp, nil
}

// Update updates an existing OTP record
func (r *otpRepository) Update(ctx context.Context, otp *models.OTPVerification) error {
	query := `
		UPDATE otp_verifications 
		SET is_verified = $1, verified_at = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, otp.IsVerified, otp.VerifiedAt, otp.ID)
	if err != nil {
		return fmt.Errorf("failed to update OTP: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("OTP not found")
	}

	return nil
}

// DeleteExpired deletes all expired OTP records
func (r *otpRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM otp_verifications 
		WHERE expires_at < $1 AND is_verified = FALSE`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired OTPs: %w", err)
	}

	return nil
}

// DeleteByUserAndPurpose deletes OTP records for a specific user and purpose
func (r *otpRepository) DeleteByUserAndPurpose(ctx context.Context, userID int, purpose string) error {
	query := `
		DELETE FROM otp_verifications 
		WHERE user_id = $1 AND purpose = $2`

	_, err := r.db.ExecContext(ctx, query, userID, purpose)
	if err != nil {
		return fmt.Errorf("failed to delete OTP: %w", err)
	}

	return nil
} 
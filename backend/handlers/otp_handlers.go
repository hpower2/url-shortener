package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpower2/url-shortener/internal/errors"
	"github.com/hpower2/url-shortener/internal/models"
	"github.com/hpower2/url-shortener/internal/repository"
	"github.com/hpower2/url-shortener/internal/services"
)

type OTPHandler struct {
	otpService         services.OTPService
	emailQueueConsumer *services.EmailQueueConsumer
	userRepo           repository.UserRepository
}

func NewOTPHandler(otpService services.OTPService, emailQueueConsumer *services.EmailQueueConsumer, userRepo repository.UserRepository) *OTPHandler {
	return &OTPHandler{
		otpService:         otpService,
		emailQueueConsumer: emailQueueConsumer,
		userRepo:           userRepo,
	}
}

// GenerateOTP generates and sends OTP to user's email
func (h *OTPHandler) GenerateOTP(c *gin.Context) {
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate OTP
	otpResponse, err := h.otpService.GenerateOTP(c.Request.Context(), user.ID, req.Email, req.Purpose)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Send OTP email via queue
	if err := h.emailQueueConsumer.PublishOTPEmail(req.Email, "", req.Purpose); err != nil {
		// Log error but don't fail the request
		// The OTP is already generated and stored
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email"})
		return
	}

	c.JSON(http.StatusOK, otpResponse)
}

// VerifyOTP verifies the provided OTP
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req models.OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify OTP
	response, err := h.otpService.VerifyOTP(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// If email verification was successful, send welcome email
	if response.IsVerified && req.Purpose == "email_verification" {
		user, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email)
		if err == nil {
			// Send welcome email via queue (non-blocking)
			if err := h.emailQueueConsumer.PublishWelcomeEmail(req.Email, user.FirstName); err != nil {
				// Log error but don't fail the response
				// The verification was successful
			}
		}
	}

	c.JSON(http.StatusOK, response)
}

// handleError handles different types of errors
func (h *OTPHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	// Handle standard errors
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}

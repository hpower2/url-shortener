package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hpower2/url-shortener/internal/config"
)

// EmailQueueConsumer handles email queue consumption and processing
type EmailQueueConsumer struct {
	rabbitMQService RabbitMQService
	emailService    EmailService
	otpService      OTPService
	config          *config.Config
}

// NewEmailQueueConsumer creates a new email queue consumer
func NewEmailQueueConsumer(
	rabbitMQService RabbitMQService,
	emailService EmailService,
	otpService OTPService,
	config *config.Config,
) *EmailQueueConsumer {
	return &EmailQueueConsumer{
		rabbitMQService: rabbitMQService,
		emailService:    emailService,
		otpService:      otpService,
		config:          config,
	}
}

// Start starts the email queue consumer
func (c *EmailQueueConsumer) Start(ctx context.Context) error {
	log.Println("Starting email queue consumer...")

	// Connect to RabbitMQ
	if err := c.rabbitMQService.Connect(); err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Start consuming emails
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Email queue consumer stopping...")
				c.rabbitMQService.Close()
				return
			default:
				if err := c.rabbitMQService.ConsumeEmails(c.handleEmailMessage); err != nil {
					log.Printf("Error consuming emails: %v", err)
					time.Sleep(5 * time.Second) // Wait before retrying
				}
			}
		}
	}()

	// Start periodic cleanup of expired OTPs
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.otpService.CleanupExpiredOTPs(context.Background()); err != nil {
					log.Printf("Error cleaning up expired OTPs: %v", err)
				}
			}
		}
	}()

	return nil
}

// handleEmailMessage processes an email message from the queue
func (c *EmailQueueConsumer) handleEmailMessage(message *EmailMessage) error {
	log.Printf("Processing email message: type=%s, to=%s", message.Type, message.To)

	switch message.Type {
	case "otp":
		return c.emailService.SendOTPEmail(message.To, message.OTPCode, message.Purpose)
	case "welcome":
		// Extract first name from the message or use a default
		firstName := "User" // You might want to pass this in the message
		return c.emailService.SendWelcomeEmail(message.To, firstName)
	default:
		return fmt.Errorf("unknown email type: %s", message.Type)
	}
}

// PublishOTPEmail publishes an OTP email to the queue
func (c *EmailQueueConsumer) PublishOTPEmail(email, otpCode, purpose string) error {
	message := &EmailMessage{
		To:         email,
		Type:       "otp",
		OTPCode:    otpCode,
		Purpose:    purpose,
		Retry:      0,
		MaxRetries: 3,
	}

	return c.rabbitMQService.PublishEmail(message)
}

// PublishWelcomeEmail publishes a welcome email to the queue
func (c *EmailQueueConsumer) PublishWelcomeEmail(email, firstName string) error {
	message := &EmailMessage{
		To:         email,
		Type:       "welcome",
		Retry:      0,
		MaxRetries: 3,
	}

	return c.rabbitMQService.PublishEmail(message)
}

// Stop stops the email queue consumer
func (c *EmailQueueConsumer) Stop() error {
	return c.rabbitMQService.Close()
}

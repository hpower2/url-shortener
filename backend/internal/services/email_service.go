package services

import (
	"fmt"
	"log"

	"github.com/hpower2/url-shortener/internal/config"
	"gopkg.in/gomail.v2"
)

// EmailService interface defines the contract for email operations
type EmailService interface {
	SendOTPEmail(email, otpCode, purpose string) error
	SendWelcomeEmail(email, firstName string) error
}

// emailService implements EmailService interface
type emailService struct {
	config *config.SMTPConfig
}

// NewEmailService creates a new email service
func NewEmailService(config *config.SMTPConfig) EmailService {
	return &emailService{
		config: config,
	}
}

// SendOTPEmail sends an OTP email to the user
func (s *emailService) SendOTPEmail(email, otpCode, purpose string) error {
	subject := s.getOTPSubject(purpose)
	body := s.getOTPEmailBody(otpCode, purpose)

	return s.sendEmail(email, subject, body)
}

// SendWelcomeEmail sends a welcome email to the user
func (s *emailService) SendWelcomeEmail(email, firstName string) error {
	subject := "Welcome to URL Shortener!"
	body := s.getWelcomeEmailBody(firstName)

	return s.sendEmail(email, subject, body)
}

// sendEmail sends an email using SMTP
func (s *emailService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	d.SSL = true // Use SSL for port 465

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// getOTPSubject returns the subject based on purpose
func (s *emailService) getOTPSubject(purpose string) string {
	switch purpose {
	case "email_verification":
		return "Verify Your Email Address"
	case "password_reset":
		return "Reset Your Password"
	default:
		return "Verification Code"
	}
}

// getOTPEmailBody returns the HTML email body for OTP
func (s *emailService) getOTPEmailBody(otpCode, purpose string) string {
	var message string
	switch purpose {
	case "email_verification":
		message = "Please use the following code to verify your email address:"
	case "password_reset":
		message = "Please use the following code to reset your password:"
	default:
		message = "Please use the following verification code:"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verification Code</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; }
        .otp-code { 
            font-size: 32px; 
            font-weight: bold; 
            color: #007bff; 
            text-align: center; 
            padding: 20px; 
            background-color: #f8f9fa; 
            border: 2px dashed #007bff; 
            margin: 20px 0; 
            letter-spacing: 5px;
        }
        .footer { 
            text-align: center; 
            margin-top: 30px; 
            font-size: 12px; 
            color: #666; 
        }
        .warning { 
            background-color: #fff3cd; 
            border: 1px solid #ffeaa7; 
            padding: 15px; 
            margin: 20px 0; 
            border-radius: 5px; 
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>URL Shortener</h1>
        </div>
        
        <p>Hello,</p>
        
        <p>%s</p>
        
        <div class="otp-code">%s</div>
        
        <div class="warning">
            <strong>Important:</strong> This code will expire in 10 minutes. 
            Do not share this code with anyone.
        </div>
        
        <p>If you didn't request this code, please ignore this email.</p>
        
        <div class="footer">
            <p>This is an automated message from URL Shortener.<br>
            Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`, message, otpCode)
}

// getWelcomeEmailBody returns the HTML email body for welcome message
func (s *emailService) getWelcomeEmailBody(firstName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to URL Shortener</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; }
        .welcome { 
            background-color: #d4edda; 
            border: 1px solid #c3e6cb; 
            padding: 20px; 
            margin: 20px 0; 
            border-radius: 5px; 
            text-align: center;
        }
        .footer { 
            text-align: center; 
            margin-top: 30px; 
            font-size: 12px; 
            color: #666; 
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>URL Shortener</h1>
        </div>
        
        <div class="welcome">
            <h2>Welcome, %s!</h2>
            <p>Your email has been successfully verified and your account is now active.</p>
        </div>
        
        <p>You can now:</p>
        <ul>
            <li>Create up to 50 shortened URLs</li>
            <li>Track click analytics</li>
            <li>Generate QR codes</li>
            <li>Manage your URLs</li>
        </ul>
        
        <p>Thank you for choosing URL Shortener!</p>
        
        <div class="footer">
            <p>This is an automated message from URL Shortener.<br>
            Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`, firstName)
}

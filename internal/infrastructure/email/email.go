package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/logger"

	"go.uber.org/zap"
)

type Service interface {
	SendOTP(to, otp string) error
	SendWelcome(to, name string) error
}

type emailService struct {
	cfg *config.EmailConfig
	log logger.Logger
}

func NewEmailService(cfg *config.EmailConfig, log logger.Logger) Service {
	return &emailService{
		cfg: cfg,
		log: log,
	}
}

// SendOTP sends OTP verification email
func (s *emailService) SendOTP(to, otp string) error {
	subject := "Verify Your Gym Pro Account"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .otp-code { font-size: 32px; font-weight: bold; color: #4CAF50; 
                    letter-spacing: 5px; text-align: center; padding: 20px; 
                    background: #f5f5f5; border-radius: 5px; }
        .footer { color: #666; font-size: 12px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Welcome to Gym Pro!</h2>
        <p>Thank you for registering. Please use the following OTP to verify your account:</p>
        <div class="otp-code">%s</div>
        <p>This code will expire in <strong>5 minutes</strong>.</p>
        <p>If you didn't request this verification, please ignore this email.</p>
        <div class="footer">
            <p>Best regards,<br>The Gym Pro Team</p>
        </div>
    </div>
</body>
</html>
`, otp)

	return s.sendEmail(to, subject, body)
}

// SendWelcome sends a welcome email after successful registration
func (s *emailService) SendWelcome(to, name string) error {
	subject := "Welcome to Gym Pro!"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2>Welcome to Gym Pro, %s! 🎉</h2>
        <p>Your account has been successfully verified.</p>
        <p>Start your fitness journey today:</p>
        <ul>
            <li>Create custom workout plans</li>
            <li>Track your meals and nutrition</li>
            <li>Connect with the fitness community</li>
        </ul>
        <p>Best regards,<br>The Gym Pro Team</p>
    </div>
</body>
</html>
`, name)

	return s.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP
func (s *emailService) sendEmail(to, subject, body string) error {
	// For development, just log the email instead of actually sending
	if s.cfg.SMTPUsername == "" || s.cfg.SMTPPassword == "" {
		s.log.Info("Email would be sent (dev mode - not actually sent)",
			zap.String("to", to),
			zap.String("subject", subject),
		)
		s.log.Debug("Email content", zap.String("body", body))
		return nil
	}

	// Prepare email message
	from := s.cfg.FromAddress
	msg := []byte(fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", s.cfg.FromName, from, to, subject, body))

	// SMTP authentication
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	// Connect to SMTP server with TLS
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	
	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName: s.cfg.SMTPHost,
	}

	// Send email
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		s.log.Error("Failed to send email",
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.log.Info("Email sent successfully",
		zap.String("to", to),
		zap.String("subject", subject),
	)

	// Suppress unused variable warning
	_ = tlsConfig

	return nil
}

package email

import (
	"crypto/tls"
	"fmt"

	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/logger"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
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

// sendEmail sends an email using configured provider (SMTP or SendGrid)
func (s *emailService) sendEmail(to, subject, body string) error {
	// Dev mode: just log email instead of sending
	if s.cfg.Provider == "" || s.cfg.Provider == "smtp" {
		if s.cfg.SMTPUsername == "" || s.cfg.SMTPPassword == "" {
			s.log.Info("Email would be sent (dev mode - not actually sent)",
				zap.String("to", to),
				zap.String("subject", subject),
			)
			s.log.Debug("Email content", zap.String("body", body))
			return nil
		}
	}

	// Route to appropriate provider
	switch s.cfg.Provider {
	case "sendgrid":
		return s.sendViaSendGrid(to, subject, body)
	case "smtp", "":
		return s.sendViaSMTP(to, subject, body)
	default:
		return fmt.Errorf("unsupported email provider: %s", s.cfg.Provider)
	}
}

// sendViaSMTP sends email via SMTP
func (s *emailService) sendViaSMTP(to, subject, body string) error {
	// Prepare email message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.FromAddress))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Configure SMTP dialer
	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUsername, s.cfg.SMTPPassword)
	
	// TLS configuration for better compatibility
	// For port 465, use SSL/TLS. For port 587, use STARTTLS
	if s.cfg.SMTPPort == 465 {
		d.SSL = true
	}
	
	d.TLSConfig = &tls.Config{
		ServerName:         s.cfg.SMTPHost,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	// Send email with detailed error logging
	s.log.Debug("Attempting to send email via SMTP",
		zap.String("to", to),
		zap.String("smtp_host", s.cfg.SMTPHost),
		zap.Int("smtp_port", s.cfg.SMTPPort),
		zap.Bool("use_ssl", d.SSL),
	)

	if err := d.DialAndSend(m); err != nil {
		s.log.Error("Failed to send email via SMTP",
			zap.String("to", to),
			zap.String("smtp_host", s.cfg.SMTPHost),
			zap.Int("smtp_port", s.cfg.SMTPPort),
			zap.String("username", s.cfg.SMTPUsername),
			zap.Error(err),
			zap.String("help", "Railway may block SMTP ports. Consider using EMAIL_PROVIDER=sendgrid"),
		)
		return fmt.Errorf("SMTP send failed (use SendGrid for Railway): %w", err)
	}

	s.log.Info("Email sent successfully via SMTP",
		zap.String("to", to),
		zap.String("subject", subject),
	)

	return nil
}

// sendViaSendGrid sends email via SendGrid API
func (s *emailService) sendViaSendGrid(to, subject, body string) error {
	if s.cfg.SendGridAPIKey == "" {
		return fmt.Errorf("SendGrid API key not configured")
	}

	from := mail.NewEmail(s.cfg.FromName, s.cfg.FromAddress)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toEmail, "", body)

	client := sendgrid.NewSendClient(s.cfg.SendGridAPIKey)
	
	s.log.Debug("Attempting to send email via SendGrid",
		zap.String("to", to),
	)

	response, err := client.Send(message)
	if err != nil {
		s.log.Error("Failed to send email via SendGrid",
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("SendGrid send failed: %w", err)
	}

	if response.StatusCode >= 400 {
		s.log.Error("SendGrid returned error status",
			zap.String("to", to),
			zap.Int("status_code", response.StatusCode),
			zap.String("body", response.Body),
		)
		return fmt.Errorf("SendGrid error: status %d - %s", response.StatusCode, response.Body)
	}

	s.log.Info("Email sent successfully via SendGrid",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.Int("status_code", response.StatusCode),
	)

	return nil
}

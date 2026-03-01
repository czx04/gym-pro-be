package email

import (
	"fmt"

	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/infrastructure/logger"

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

// sendEmail sends an email using SMTP
func (s *emailService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.FromAddress)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUsername, s.cfg.SMTPPassword)

	// Nếu dùng cổng 465 thì để mặc định, nếu 587 thư viện tự hiểu STARTTLS
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

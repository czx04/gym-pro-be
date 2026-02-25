package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost
	DefaultCost = bcrypt.DefaultCost
)

// PasswordManager handles password hashing and verification
type PasswordManager struct {
	cost int
}

// NewPasswordManager creates a new password manager
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{
		cost: DefaultCost,
	}
}

// HashPassword hashes a password using bcrypt
func (p *PasswordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a hashed password with a plaintext password
func (p *PasswordManager) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// VerifyPassword checks if a password matches the hash
func (p *PasswordManager) VerifyPassword(hashedPassword, password string) bool {
	err := p.ComparePassword(hashedPassword, password)
	return err == nil
}

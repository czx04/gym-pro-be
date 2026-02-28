package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost = bcrypt.DefaultCost
)

type PasswordManager struct {
	cost int
}

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{
		cost: DefaultCost,
	}
}

func (p *PasswordManager) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (p *PasswordManager) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (p *PasswordManager) VerifyPassword(hashedPassword, password string) bool {
	err := p.ComparePassword(hashedPassword, password)
	return err == nil
}

package auth

import (
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret             string
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
}

func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		secret:             cfg.Secret,
		accessTokenExpire:  cfg.AccessTokenExpire,
		refreshTokenExpire: cfg.RefreshTokenExpire,
	}
}

func (j *JWTManager) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.accessTokenExpire)

	claims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token
func (j *JWTManager) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshTokenExpire)

	claims := Claims{
		UserID:    userID,
		Email:     email,
		TokenType: RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTManager) GenerateTokenPair(userID uuid.UUID, email string) (accessToken, refreshToken string, err error) {
	accessToken, err = j.GenerateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.GenerateRefreshToken(userID, email)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, errors.TokenInvalid()
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.TokenInvalid()
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.TokenExpired()
	}

	return claims, nil
}

// ValidateAccessToken validates an access token
func (j *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != AccessToken {
		return nil, errors.TokenInvalid()
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (j *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshToken {
		return nil, errors.TokenInvalid()
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (j *JWTManager) RefreshAccessToken(refreshTokenString string) (string, error) {
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	return j.GenerateAccessToken(claims.UserID, claims.Email)
}
